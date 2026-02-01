package matches

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context, query ListQuery) ([]MatchRow, int, error) {
	baseQuery := `FROM matched_certificates mc JOIN keywords k ON mc.keyword_id = k.id WHERE k.is_deleted = FALSE`
	args := []interface{}{}
	argIndex := 1

	if query.Keyword != "" {
		baseQuery += fmt.Sprintf(" AND k.normalized_value = $%d", argIndex)
		args = append(args, strings.ToLower(query.Keyword))
		argIndex++
	}

	if query.Q != "" {
		baseQuery += fmt.Sprintf(" AND (mc.domain_name ILIKE $%d OR mc.issuer_cn ILIKE $%d OR mc.issuer_org ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+query.Q+"%")
		argIndex++
	}

	if query.Issuer != "" {
		baseQuery += fmt.Sprintf(" AND (mc.issuer_cn ILIKE $%d OR mc.issuer_org ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+query.Issuer+"%")
		argIndex++
	}

	if query.DateFrom != nil {
		baseQuery += fmt.Sprintf(" AND mc.first_seen_at >= $%d", argIndex)
		args = append(args, *query.DateFrom)
		argIndex++
	}

	if query.DateTo != nil {
		baseQuery += fmt.Sprintf(" AND mc.first_seen_at <= $%d", argIndex)
		args = append(args, *query.DateTo)
		argIndex++
	}

	if query.NewOnly {
		baseQuery += " AND mc.is_new = TRUE"
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := `SELECT mc.id, mc.keyword_id, k.keyword, mc.certificate_sha256, mc.matched_field, mc.matched_value,
		mc.domain_name, mc.issuer_cn, mc.issuer_org, mc.subject_cn, mc.subject_org,
		mc.not_before, mc.not_after, mc.first_seen_at, mc.last_seen_at, mc.is_new, mc.ct_log_index ` + baseQuery

	switch query.Sort {
	case SortLastSeenDesc:
		selectQuery += " ORDER BY mc.last_seen_at DESC"
	case SortDomainAsc:
		selectQuery += " ORDER BY mc.domain_name ASC NULLS LAST"
	default:
		selectQuery += " ORDER BY mc.first_seen_at DESC"
	}

	offset := (query.Page - 1) * query.PageSize
	selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, query.PageSize, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matches []MatchRow
	for rows.Next() {
		var m MatchRow
		if err := rows.Scan(&m.ID, &m.KeywordID, &m.KeywordValue, &m.CertSHA256, &m.MatchedField, &m.MatchedValue,
			&m.DomainName, &m.IssuerCN, &m.IssuerOrg, &m.SubjectCN, &m.SubjectOrg,
			&m.NotBefore, &m.NotAfter, &m.FirstSeenAt, &m.LastSeenAt, &m.IsNew, &m.CtLogIndex); err != nil {
			return nil, 0, err
		}
		matches = append(matches, m)
	}

	return matches, total, nil
}

func (r *Repository) StreamAll(ctx context.Context, query ListQuery, fn func(MatchRow) error) error {
	baseQuery := `FROM matched_certificates mc JOIN keywords k ON mc.keyword_id = k.id WHERE k.is_deleted = FALSE`
	args := []interface{}{}
	argIndex := 1

	if query.Keyword != "" {
		baseQuery += fmt.Sprintf(" AND k.normalized_value = $%d", argIndex)
		args = append(args, strings.ToLower(query.Keyword))
		argIndex++
	}

	if query.Q != "" {
		baseQuery += fmt.Sprintf(" AND (mc.domain_name ILIKE $%d OR mc.issuer_cn ILIKE $%d OR mc.issuer_org ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+query.Q+"%")
		argIndex++
	}

	if query.Issuer != "" {
		baseQuery += fmt.Sprintf(" AND (mc.issuer_cn ILIKE $%d OR mc.issuer_org ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+query.Issuer+"%")
		argIndex++
	}

	if query.DateFrom != nil {
		baseQuery += fmt.Sprintf(" AND mc.first_seen_at >= $%d", argIndex)
		args = append(args, *query.DateFrom)
		argIndex++
	}

	if query.DateTo != nil {
		baseQuery += fmt.Sprintf(" AND mc.first_seen_at <= $%d", argIndex)
		args = append(args, *query.DateTo)
		argIndex++
	}

	if query.NewOnly {
		baseQuery += " AND mc.is_new = TRUE"
	}

	selectQuery := `SELECT mc.id, mc.keyword_id, k.keyword, mc.certificate_sha256, mc.matched_field, mc.matched_value,
		mc.domain_name, mc.issuer_cn, mc.issuer_org, mc.subject_cn, mc.subject_org,
		mc.not_before, mc.not_after, mc.first_seen_at, mc.last_seen_at, mc.is_new, mc.ct_log_index ` + baseQuery

	switch query.Sort {
	case SortLastSeenDesc:
		selectQuery += " ORDER BY mc.last_seen_at DESC"
	case SortDomainAsc:
		selectQuery += " ORDER BY mc.domain_name ASC NULLS LAST"
	default:
		selectQuery += " ORDER BY mc.first_seen_at DESC"
	}

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var m MatchRow
		if err := rows.Scan(&m.ID, &m.KeywordID, &m.KeywordValue, &m.CertSHA256, &m.MatchedField, &m.MatchedValue,
			&m.DomainName, &m.IssuerCN, &m.IssuerOrg, &m.SubjectCN, &m.SubjectOrg,
			&m.NotBefore, &m.NotAfter, &m.FirstSeenAt, &m.LastSeenAt, &m.IsNew, &m.CtLogIndex); err != nil {
			return err
		}
		if err := fn(m); err != nil {
			return err
		}
	}

	return rows.Err()
}
