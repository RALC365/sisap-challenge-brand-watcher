package scheduler

import (
	"context"
	"errors"
	"sync"
	"time"

	"brand-protection-monitor/internal/ct"
	"brand-protection-monitor/internal/db"
	"brand-protection-monitor/internal/matcher"
	"brand-protection-monitor/internal/parser"

	"go.uber.org/zap"
)

type Scheduler struct {
	logger   *zap.Logger
	repo     *Repository
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

func New(logger *zap.Logger) *Scheduler {
	return &Scheduler{
		logger:   logger,
		repo:     NewRepository(db.GetPool()),
		stopChan: make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	s.wg.Add(1)
	go s.run(ctx)

	s.logger.Info("scheduler started")
}

func (s *Scheduler) run(ctx context.Context) {
	defer s.wg.Done()

	config, err := s.repo.GetConfig(ctx)
	if err != nil {
		s.logger.Error("failed to load monitor config", zap.Error(err))
		return
	}

	pollInterval := time.Duration(config.PollIntervalSec) * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	s.runCycle(ctx)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopped due to context cancellation")
			return
		case <-s.stopChan:
			s.logger.Info("scheduler stopped")
			return
		case <-ticker.C:
			s.runCycle(ctx)
		}
	}
}

func (s *Scheduler) runCycle(ctx context.Context) {
	s.logger.Debug("starting CT log polling cycle")

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return
	}

	state, err := s.repo.LockStateForUpdate(ctx, tx)
	if err != nil {
		tx.Rollback(ctx)
		s.logger.Error("failed to lock state", zap.Error(err))
		return
	}

	if state.State != "idle" {
		tx.Rollback(ctx)
		s.logger.Debug("monitor not idle, skipping cycle", zap.String("state", state.State))
		return
	}

	if err := s.repo.SetStateRunning(ctx, tx); err != nil {
		tx.Rollback(ctx)
		s.logger.Error("failed to set state running", zap.Error(err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("failed to commit running state", zap.Error(err))
		return
	}

	config, err := s.repo.GetConfig(ctx)
	if err != nil {
		s.handleFatalError(ctx, "CONFIG_ERROR", err.Error())
		return
	}

	s.executeCycle(ctx, config, state)
}

func (s *Scheduler) executeCycle(ctx context.Context, config *MonitorConfig, state *MonitorState) {
	ctClient := ct.NewClient(ct.ClientConfig{
		BaseURL:        config.CTLogBaseURL,
		ConnectTimeout: time.Duration(config.ConnectTimeoutMs) * time.Millisecond,
		ReadTimeout:    time.Duration(config.ReadTimeoutMs) * time.Millisecond,
		BatchSize:      config.BatchSize,
	})

	sth, err := ctClient.GetSTH(ctx)
	if err != nil {
		if errors.Is(err, ct.ErrInvalidJSON) {
			s.handleFatalError(ctx, "CT_INVALID_JSON", err.Error())
		} else if errors.Is(err, ct.ErrTimeout) {
			s.handleFatalError(ctx, "CT_TIMEOUT", err.Error())
		} else {
			s.handleFatalError(ctx, "CT_CONNECTION_ERROR", err.Error())
		}
		return
	}

	s.logger.Info("got STH", zap.Int64("tree_size", sth.TreeSize))

	fetchRange := ct.CalculateRange(sth.TreeSize, config.BatchSize, state.LastProcessedIndex)
	s.logger.Info("calculated range", zap.Int64("start", fetchRange.Start), zap.Int64("end", fetchRange.End))

	if fetchRange.Start > fetchRange.End {
		s.logger.Info("no new entries to process")
		if err := s.repo.SetStateIdle(ctx, sth.TreeSize, state.LastProcessedIndex); err != nil {
			s.logger.Error("failed to set state idle", zap.Error(err))
		}
		return
	}

	runID, err := s.repo.CreateRun(ctx, fetchRange.Start, fetchRange.End)
	if err != nil {
		s.handleFatalError(ctx, "DB_ERROR", err.Error())
		return
	}

	entries, err := ctClient.GetEntriesChunked(ctx, fetchRange.Start, fetchRange.End, 100)
	if err != nil {
		if errors.Is(err, ct.ErrInvalidJSON) {
			s.handleFatalError(ctx, "CT_INVALID_JSON", err.Error())
		} else {
			s.handleFatalError(ctx, "CT_FETCH_ERROR", err.Error())
		}
		s.repo.UpdateRunError(ctx, runID, "CT_ERROR", err.Error())
		return
	}

	s.logger.Info("fetched entries", zap.Int("count", len(entries)))

	keywordRows, err := s.repo.GetActiveKeywords(ctx)
	if err != nil {
		s.handleFatalError(ctx, "DB_ERROR", err.Error())
		s.repo.UpdateRunError(ctx, runID, "DB_ERROR", err.Error())
		return
	}

	keywords := make([]matcher.Keyword, len(keywordRows))
	for i, kw := range keywordRows {
		keywords[i] = matcher.Keyword{
			ID:              kw.ID,
			Value:           kw.Keyword,
			NormalizedValue: kw.NormalizedValue,
		}
	}

	m := matcher.New(keywords)

	entriesFetched := len(entries)
	entriesParsed := 0
	parseErrors := 0
	matchesFound := 0
	lastProcessedIndex := fetchRange.Start

	for i, entry := range entries {
		index := fetchRange.Start + int64(i)
		result := parser.ParseEntry(entry, index)

		if result.Error != nil {
			parseErrors++
			s.logger.Debug("parse error", zap.Int64("index", index), zap.Error(result.Error))
			continue
		}

		entriesParsed++
		cert := result.Certificate

		matches := m.Match(cert)
		for _, match := range matches {
			matchesFound++

			insert := MatchInsert{
				KeywordID:       match.KeywordID,
				MonitorRunID:    runID,
				CertFingerprint: cert.Fingerprint,
				MatchedField:    string(match.MatchedField),
				MatchedValue:    match.MatchedValue,
				CTLogIndex:      index,
				NotBefore:       cert.NotBefore,
				NotAfter:        cert.NotAfter,
				SANList:         cert.SANs,
			}

			if cert.SubjectCN != "" {
				insert.SubjectCN = &cert.SubjectCN
			}
			if cert.SubjectOrg != "" {
				insert.SubjectOrg = &cert.SubjectOrg
			}
			if cert.IssuerCN != "" {
				insert.IssuerCN = &cert.IssuerCN
			}
			if cert.IssuerOrg != "" {
				insert.IssuerOrg = &cert.IssuerOrg
			}
			if match.DomainName != "" {
				insert.DomainName = &match.DomainName
			}
			ctLogURL := config.CTLogBaseURL
			insert.CTLogURL = &ctLogURL

			if err := s.repo.UpsertMatch(ctx, insert); err != nil {
				s.logger.Error("failed to upsert match", zap.Error(err))
			}
		}

		lastProcessedIndex = index
	}

	if err := s.repo.UpdateRunSuccess(ctx, runID, entriesFetched, entriesParsed, parseErrors, matchesFound, lastProcessedIndex); err != nil {
		s.logger.Error("failed to update run", zap.Error(err))
	}

	if err := s.repo.SetStateIdle(ctx, sth.TreeSize, lastProcessedIndex); err != nil {
		s.logger.Error("failed to set state idle", zap.Error(err))
	}

	s.logger.Info("cycle completed",
		zap.Int("entries_fetched", entriesFetched),
		zap.Int("entries_parsed", entriesParsed),
		zap.Int("parse_errors", parseErrors),
		zap.Int("matches_found", matchesFound))
}

func (s *Scheduler) handleFatalError(ctx context.Context, code, message string) {
	s.logger.Error("fatal cycle error", zap.String("code", code), zap.String("message", message))
	if err := s.repo.SetStateError(ctx, code, message); err != nil {
		s.logger.Error("failed to set error state", zap.Error(err))
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	close(s.stopChan)
	s.wg.Wait()
	s.running = false
}
