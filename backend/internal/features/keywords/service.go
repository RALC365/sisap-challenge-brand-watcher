package keywords

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, query ListQuery) (*ListKeywordsResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize != 10 && query.PageSize != 25 && query.PageSize != 50 {
		query.PageSize = 10
	}

	rows, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]KeywordItem, len(rows))
	for i, row := range rows {
		items[i] = KeywordItem{
			KeywordID:       row.ID,
			Value:           row.Keyword,
			NormalizedValue: row.NormalizedValue,
			Status:          row.Status,
			CreatedAt:       row.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
		}
	}

	return &ListKeywordsResponse{
		Items: items,
		Total: total,
	}, nil
}

func (s *Service) Create(ctx context.Context, req CreateKeywordRequest) (*Keyword, error) {
	value := strings.TrimSpace(req.Value)
	if value == "" {
		return nil, ErrEmptyValue
	}
	if len(value) < 1 || len(value) > 64 {
		return nil, ErrValueTooLong
	}

	normalizedValue := NormalizeValue(value)

	row, err := s.repo.Create(ctx, value, normalizedValue)
	if err != nil {
		return nil, err
	}

	return &Keyword{
		ID:              row.ID,
		Value:           row.Keyword,
		NormalizedValue: row.NormalizedValue,
		Status:          KeywordStatus(row.Status),
		CreatedAt:       row.CreatedAt,
	}, nil
}

func (s *Service) Delete(ctx context.Context, keywordID string) error {
	if _, err := uuid.Parse(keywordID); err != nil {
		return ErrInvalidKeywordID
	}

	rowsAffected, err := s.repo.SoftDelete(ctx, keywordID)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrKeywordNotFound
	}

	return nil
}
