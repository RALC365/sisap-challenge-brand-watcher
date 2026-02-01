package export

import (
	"context"
	"fmt"
	"io"
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

	header := "id,keyword_id,keyword_value,certificate_sha256,matched_field,matched_value,domain_name,issuer_cn,issuer_org,subject_cn,subject_org,not_before,not_after,first_seen_at,last_seen_at,is_new,ct_log_index\n"
	if _, err := writer.Write([]byte(header)); err != nil {
		s.repo.UpdateFailed(ctx, exportID, err.Error())
		return 0, err
	}

	count := 0
	err = s.matchService.StreamAll(ctx, query, func(m matches.Match) error {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%t,%d\n",
			escapeCSV(m.ID),
			escapeCSV(m.KeywordID),
			escapeCSV(m.KeywordValue),
			escapeCSV(m.CertSHA256),
			escapeCSV(m.MatchedField),
			escapeCSV(m.MatchedValue),
			escapeCSVPtr(m.DomainName),
			escapeCSVPtr(m.IssuerCN),
			escapeCSVPtr(m.IssuerOrg),
			escapeCSVPtr(m.SubjectCN),
			escapeCSVPtr(m.SubjectOrg),
			formatTimePtr(m.NotBefore),
			formatTimePtr(m.NotAfter),
			m.FirstSeenAt.Format(time.RFC3339),
			m.LastSeenAt.Format(time.RFC3339),
			m.IsNew,
			m.CtLogIndex,
		)
		if _, err := writer.Write([]byte(line)); err != nil {
			return err
		}
		count++
		return nil
	})

	if err != nil {
		s.repo.UpdateFailed(ctx, exportID, err.Error())
		return count, err
	}

	if err := s.repo.UpdateCompleted(ctx, exportID, count); err != nil {
		return count, nil
	}

	return count, nil
}

func escapeCSV(s string) string {
	if s == "" {
		return ""
	}
	needsQuotes := false
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' || c == '\r' {
			needsQuotes = true
			break
		}
	}
	if !needsQuotes {
		return s
	}
	escaped := ""
	for _, c := range s {
		if c == '"' {
			escaped += "\"\""
		} else {
			escaped += string(c)
		}
	}
	return "\"" + escaped + "\""
}

func escapeCSVPtr(s *string) string {
	if s == nil {
		return ""
	}
	return escapeCSV(*s)
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
