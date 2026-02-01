package scheduler

import (
        "context"
        "time"

        "github.com/jackc/pgx/v5"
        "github.com/jackc/pgx/v5/pgxpool"
)

type MonitorConfig struct {
        CTLogBaseURL      string
        PollIntervalSec   int
        BatchSize         int
        ConnectTimeoutMs  int
        ReadTimeoutMs     int
}

type MonitorState struct {
        State              string
        LastTreeSize       int64
        LastProcessedIndex int64
        LastRunAt          *time.Time
        LastSuccessAt      *time.Time
        LastErrorCode      *string
        LastError          *string
}

type MonitorRun struct {
        ID              string
        StartedAt       time.Time
        FinishedAt      *time.Time
        Status          string
        StartIndex      int64
        EndIndex        int64
        EntriesFetched  int
        EntriesParsed   int
        ParseErrorCount int
        MatchesFound    int
        ErrorCode       *string
        ErrorMessage    *string
}

type Repository struct {
        pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
        return &Repository{pool: pool}
}

func (r *Repository) GetConfig(ctx context.Context) (*MonitorConfig, error) {
        var cfg MonitorConfig
        err := r.pool.QueryRow(ctx, `
                SELECT ct_log_base_url, poll_interval_seconds, batch_size, ct_connect_timeout_ms, ct_read_timeout_ms
                FROM monitor_config WHERE id = 1
        `).Scan(&cfg.CTLogBaseURL, &cfg.PollIntervalSec, &cfg.BatchSize, &cfg.ConnectTimeoutMs, &cfg.ReadTimeoutMs)
        if err != nil {
                return nil, err
        }
        return &cfg, nil
}

func (r *Repository) GetState(ctx context.Context) (*MonitorState, error) {
        var state MonitorState
        err := r.pool.QueryRow(ctx, `
                SELECT state, last_tree_size, last_processed_index, last_run_at, last_success_at, last_error_code, last_error
                FROM monitor_state WHERE id = 1
        `).Scan(&state.State, &state.LastTreeSize, &state.LastProcessedIndex, &state.LastRunAt, &state.LastSuccessAt, &state.LastErrorCode, &state.LastError)
        if err != nil {
                return nil, err
        }
        return &state, nil
}

func (r *Repository) LockStateForUpdate(ctx context.Context, tx pgx.Tx) (*MonitorState, error) {
        var state MonitorState
        err := tx.QueryRow(ctx, `
                SELECT state, last_tree_size, last_processed_index, last_run_at, last_success_at, last_error_code, last_error
                FROM monitor_state WHERE id = 1 FOR UPDATE
        `).Scan(&state.State, &state.LastTreeSize, &state.LastProcessedIndex, &state.LastRunAt, &state.LastSuccessAt, &state.LastErrorCode, &state.LastError)
        if err != nil {
                return nil, err
        }
        return &state, nil
}

func (r *Repository) SetStateRunning(ctx context.Context, tx pgx.Tx) error {
        _, err := tx.Exec(ctx, `
                UPDATE monitor_state SET state = 'running', last_run_at = NOW(), last_error_code = NULL, last_error = NULL WHERE id = 1
        `)
        return err
}

func (r *Repository) SetStateIdle(ctx context.Context, treeSize, processedIndex int64) error {
        _, err := r.pool.Exec(ctx, `
                UPDATE monitor_state SET state = 'idle', last_tree_size = $1, last_processed_index = $2, last_success_at = NOW() WHERE id = 1
        `, treeSize, processedIndex)
        return err
}

func (r *Repository) SetStateError(ctx context.Context, errorCode, errorMessage string) error {
        _, err := r.pool.Exec(ctx, `
                UPDATE monitor_state SET state = 'error', last_error_code = $1, last_error = $2 WHERE id = 1
        `, errorCode, errorMessage)
        return err
}

func (r *Repository) CreateRun(ctx context.Context, startIndex, endIndex int64) (string, error) {
        var id string
        err := r.pool.QueryRow(ctx, `
                INSERT INTO monitor_run (start_index, end_index, state)
                VALUES ($1, $2, 'running')
                RETURNING id
        `, startIndex, endIndex).Scan(&id)
        return id, err
}

func (r *Repository) UpdateRunSuccess(ctx context.Context, runID string, entriesFetched, entriesParsed, parseErrors, matchesFound int, endIndex int64) error {
        _, err := r.pool.Exec(ctx, `
                UPDATE monitor_run SET 
                        ended_at = NOW(), 
                        state = 'completed',
                        certificates_processed = $2,
                        parse_error_count = $3,
                        matches_found = $4,
                        end_index = $5
                WHERE id = $1
        `, runID, entriesParsed, parseErrors, matchesFound, endIndex)
        return err
}

func (r *Repository) UpdateRunError(ctx context.Context, runID string, errorCode, errorMessage string) error {
        _, err := r.pool.Exec(ctx, `
                UPDATE monitor_run SET 
                        ended_at = NOW(), 
                        state = 'failed',
                        error_message = $2
                WHERE id = $1
        `, runID, errorMessage)
        return err
}

func (r *Repository) GetActiveKeywords(ctx context.Context) ([]KeywordRow, error) {
        rows, err := r.pool.Query(ctx, `
                SELECT id, keyword, normalized_value FROM keywords WHERE is_deleted = FALSE AND status = 'active'
        `)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var keywords []KeywordRow
        for rows.Next() {
                var kw KeywordRow
                if err := rows.Scan(&kw.ID, &kw.Keyword, &kw.NormalizedValue); err != nil {
                        return nil, err
                }
                keywords = append(keywords, kw)
        }
        return keywords, nil
}

type KeywordRow struct {
        ID              string
        Keyword         string
        NormalizedValue string
}

func (r *Repository) UpsertMatch(ctx context.Context, m MatchInsert) error {
        _, err := r.pool.Exec(ctx, `
                INSERT INTO matched_certificates (
                        keyword_id, monitor_run_id, certificate_sha256, matched_field, matched_value,
                        domain_name, issuer_cn, issuer_org, subject_cn, subject_org, san_list,
                        not_before, not_after, ct_log_index, ct_log_url, first_seen_at, last_seen_at, is_new
                ) VALUES (
                        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW(), TRUE
                )
                ON CONFLICT (certificate_sha256, keyword_id, matched_field) 
                DO UPDATE SET 
                        last_seen_at = NOW(),
                        monitor_run_id = EXCLUDED.monitor_run_id,
                        is_new = FALSE
        `, m.KeywordID, m.MonitorRunID, m.CertFingerprint, m.MatchedField, m.MatchedValue,
                m.DomainName, m.IssuerCN, m.IssuerOrg, m.SubjectCN, m.SubjectOrg, m.SANList,
                m.NotBefore, m.NotAfter, m.CTLogIndex, m.CTLogURL)
        return err
}

type MatchInsert struct {
        KeywordID       string
        MonitorRunID    string
        CertFingerprint string
        MatchedField    string
        MatchedValue    string
        DomainName      *string
        IssuerCN        *string
        IssuerOrg       *string
        SubjectCN       *string
        SubjectOrg      *string
        SANList         []string
        NotBefore       *time.Time
        NotAfter        *time.Time
        CTLogIndex      int64
        CTLogURL        *string
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
        return r.pool.Begin(ctx)
}
