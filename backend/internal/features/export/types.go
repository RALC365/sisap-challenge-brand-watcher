package export

import (
	"errors"
	"time"

	"brand-protection-monitor/internal/features/matches"
)

type ExportStatus string

const (
	ExportStatusPending   ExportStatus = "pending"
	ExportStatusStreaming ExportStatus = "streaming"
	ExportStatusCompleted ExportStatus = "completed"
	ExportStatusFailed    ExportStatus = "failed"
)

var (
	ErrRateLimited = errors.New("rate limit exceeded")
)

const (
	ErrorCodeRateLimited = "RATE_LIMITED"
	ErrorCodeExportError = "EXPORT_ERROR"
	ErrorCodeDBError     = "DB_ERROR"
)

type ExportRecord struct {
	ID           string
	Filename     string
	RecordCount  int
	FilterParams FilterParams
	Status       ExportStatus
	ErrorMessage *string
	CreatedAt    time.Time
}

type FilterParams struct {
	Keyword  string     `json:"keyword,omitempty"`
	Q        string     `json:"q,omitempty"`
	Issuer   string     `json:"issuer,omitempty"`
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	NewOnly  bool       `json:"new_only,omitempty"`
	Sort     string     `json:"sort,omitempty"`
}

func FiltersFromMatchQuery(q matches.ListQuery) FilterParams {
	return FilterParams{
		Keyword:  q.Keyword,
		Q:        q.Q,
		Issuer:   q.Issuer,
		DateFrom: q.DateFrom,
		DateTo:   q.DateTo,
		NewOnly:  q.NewOnly,
		Sort:     string(q.Sort),
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
