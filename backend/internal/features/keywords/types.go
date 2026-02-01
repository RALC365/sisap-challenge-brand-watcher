package keywords

import (
	"errors"
	"time"
)

type KeywordStatus string

const (
	KeywordStatusActive   KeywordStatus = "active"
	KeywordStatusInactive KeywordStatus = "inactive"
)

var (
	ErrMissingValue     = errors.New("value is required")
	ErrEmptyValue       = errors.New("value cannot be empty after trimming")
	ErrValueTooLong     = errors.New("value must be between 1 and 64 characters")
	ErrInvalidKeywordID = errors.New("keyword_id must be a valid UUID")
	ErrKeywordNotFound  = errors.New("keyword not found")
	ErrDuplicateKeyword = errors.New("keyword already exists")
)

const (
	ErrorCodeValidation      = "VALIDATION_ERROR"
	ErrorCodeDuplicate       = "DUPLICATE_KEYWORD"
	ErrorCodeNotFound        = "NOT_FOUND"
	ErrorCodeInvalidPathParam = "INVALID_PATH_PARAM"
	ErrorCodeDBError         = "DB_ERROR"
)

type Keyword struct {
	ID              string        `json:"keyword_id"`
	Value           string        `json:"value"`
	NormalizedValue string        `json:"normalized_value"`
	Status          KeywordStatus `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
}

type KeywordRow struct {
	ID              string
	Keyword         string
	NormalizedValue string
	Status          string
	IsDeleted       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ListQuery struct {
	Q        string
	Page     int
	PageSize int
}
