package keywords

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrMissingValue    = errors.New("value is required")
	ErrEmptyValue      = errors.New("value cannot be empty")
	ErrValueTooLong    = errors.New("value must be between 1 and 64 characters")
	ErrKeywordNotFound = errors.New("keyword not found")
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

	rows, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]Keyword, len(rows))
	for i, row := range rows {
		items[i] = Keyword{
			ID:              row.ID,
			Value:           row.Keyword,
			NormalizedValue: row.NormalizedValue,
			Status:          KeywordStatus(row.Status),
			CreatedAt:       row.CreatedAt,
		}
	}

	return &ListResponse{
		Items: items,
		Total: total,
	}, nil
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Keyword, error) {
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
		if errors.Is(err, ErrDuplicateKeyword) {
			return nil, ErrDuplicateKeyword
		}
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

func (s *Service) Delete(ctx context.Context, id string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrKeywordNotFound
	}

	return s.repo.SoftDelete(ctx, id)
}
