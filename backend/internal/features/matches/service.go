package matches

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, query ListQuery) (*ListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize != 10 && query.PageSize != 25 && query.PageSize != 50 {
		query.PageSize = 10
	}
	if query.Sort == "" {
		query.Sort = SortFirstSeenDesc
	}

	if query.NewOnly {
		lastSuccessTime, err := s.repo.GetLastSuccessTime(ctx)
		if err == nil && lastSuccessTime != nil {
			query.LastSuccessTime = lastSuccessTime
		}
	}

	rows, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]Match, len(rows))
	for i, row := range rows {
		items[i] = rowToMatch(row)
	}

	return &ListResponse{
		Items: items,
		Total: total,
	}, nil
}

func (s *Service) StreamAll(ctx context.Context, query ListQuery, fn func(Match) error) error {
	if query.NewOnly {
		lastSuccessTime, err := s.repo.GetLastSuccessTime(ctx)
		if err == nil && lastSuccessTime != nil {
			query.LastSuccessTime = lastSuccessTime
		}
	}

	return s.repo.StreamAll(ctx, query, func(row MatchRow) error {
		return fn(rowToMatch(row))
	})
}

func rowToMatch(row MatchRow) Match {
	return Match{
		ID:           row.ID,
		KeywordID:    row.KeywordID,
		KeywordValue: row.KeywordValue,
		CertSHA256:   row.CertSHA256,
		MatchedField: row.MatchedField,
		MatchedValue: row.MatchedValue,
		DomainName:   row.DomainName,
		IssuerCN:     row.IssuerCN,
		IssuerOrg:    row.IssuerOrg,
		SubjectCN:    row.SubjectCN,
		SubjectOrg:   row.SubjectOrg,
		NotBefore:    row.NotBefore,
		NotAfter:     row.NotAfter,
		FirstSeenAt:  row.FirstSeenAt,
		LastSeenAt:   row.LastSeenAt,
		IsNew:        row.IsNew,
		CtLogIndex:   row.CtLogIndex,
	}
}
