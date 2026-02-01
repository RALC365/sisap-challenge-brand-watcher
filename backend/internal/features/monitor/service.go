package monitor

import (
        "context"
        "errors"

        "github.com/jackc/pgx/v5"
)

type Service struct {
        repo *Repository
}

func NewService(repo *Repository) *Service {
        return &Service{repo: repo}
}

func (s *Service) GetStatus(ctx context.Context) (*StatusResponse, error) {
        state, err := s.repo.GetState(ctx)
        if err != nil {
                return nil, err
        }

        pollInterval, _ := s.repo.GetPollInterval(ctx)

        response := &StatusResponse{
                State:               MonitorState(state.State),
                LastRunAt:           state.LastRunAt,
                LastSuccessAt:       state.LastSuccessAt,
                LastErrorCode:       state.LastErrorCode,
                LastErrorMessage:    state.LastErrorMessage,
                PollIntervalSeconds: pollInterval,
        }

        lastRun, err := s.repo.GetLastCompletedRun(ctx)
        if err != nil && !errors.Is(err, pgx.ErrNoRows) {
                return nil, err
        }

        if lastRun != nil {
                metrics := &MetricsLastRun{
                        ProcessedCount:  lastRun.CertificatesProcessed,
                        MatchCount:      lastRun.MatchesFound,
                        ParseErrorCount: lastRun.ParseErrorCount,
                }
                if lastRun.DurationMs != nil {
                        metrics.DurationMs = *lastRun.DurationMs
                }
                if lastRun.CtLatencyMs != nil {
                        metrics.CtLatencyMs = *lastRun.CtLatencyMs
                }
                if lastRun.DbLatencyMs != nil {
                        metrics.DbLatencyMs = *lastRun.DbLatencyMs
                }
                response.MetricsLastRun = metrics
        }

        return response, nil
}
