package export

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"brand-protection-monitor/internal/features/matches"
)

type Service struct {
	repo         *Repository
	matchService *matches.Service
}

func NewService(repo *Repository, matchService *matches.Service) *Service {
	return &Service{
		repo:         repo,
		matchService: matchService,
	}
}

func (s *Service) ExportCSV(ctx context.Context, query matches.ListQuery, writer io.Writer) (int, error) {
	filename := fmt.Sprintf("export_%s.csv", time.Now().Format("20060102_150405"))
	filterParams := FiltersFromMatchQuery(query)

	exportID, err := s.repo.Create(ctx, filename, filterParams)
	if err != nil {
		return 0, err
	}

	if err := s.repo.UpdateStatus(ctx, exportID, ExportStatusStreaming); err != nil {
		return 0, err
	}

	csvWriter := csv.NewWriter(writer)
	csvWriter.UseCRLF = true

	header := []string{
		"id", "keyword_id", "keyword_value", "certificate_sha256",
		"matched_field", "matched_value", "domain_name",
		"issuer_cn", "issuer_org", "subject_cn", "subject_org",
		"not_before", "not_after", "first_seen_at", "last_seen_at",
		"is_new", "ct_log_index",
	}
	if err := csvWriter.Write(header); err != nil {
		s.repo.UpdateFailed(ctx, exportID, err.Error())
		return 0, err
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		s.repo.UpdateFailed(ctx, exportID, err.Error())
		return 0, err
	}

	count := 0
	err = s.matchService.StreamAll(ctx, query, func(m matches.Match) error {
		row := []string{
			m.ID,
			m.KeywordID,
			m.KeywordValue,
			m.CertSHA256,
			m.MatchedField,
			m.MatchedValue,
			ptrToString(m.DomainName),
			ptrToString(m.IssuerCN),
			ptrToString(m.IssuerOrg),
			ptrToString(m.SubjectCN),
			ptrToString(m.SubjectOrg),
			formatTimePtr(m.NotBefore),
			formatTimePtr(m.NotAfter),
			m.FirstSeenAt.Format(time.RFC3339),
			m.LastSeenAt.Format(time.RFC3339),
			strconv.FormatBool(m.IsNew),
			strconv.FormatInt(m.CtLogIndex, 10),
		}

		if err := csvWriter.Write(row); err != nil {
			return err
		}

		if count%100 == 0 {
			csvWriter.Flush()
			if err := csvWriter.Error(); err != nil {
				return err
			}
		}

		count++
		return nil
	})

	csvWriter.Flush()
	if flushErr := csvWriter.Error(); flushErr != nil && err == nil {
		err = flushErr
	}

	if err != nil {
		s.repo.UpdateFailed(ctx, exportID, err.Error())
		return count, err
	}

	if updateErr := s.repo.UpdateCompleted(ctx, exportID, count); updateErr != nil {
		return count, nil
	}

	return count, nil
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
