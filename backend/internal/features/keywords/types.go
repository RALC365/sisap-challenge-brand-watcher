package keywords

import "time"

type KeywordStatus string

const (
	KeywordStatusActive   KeywordStatus = "active"
	KeywordStatusInactive KeywordStatus = "inactive"
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

type ListResponse struct {
	Items []Keyword `json:"items"`
	Total int       `json:"total"`
}

type CreateRequest struct {
	Value string `json:"value"`
}

type DeleteResponse struct {
	OK bool `json:"ok"`
}
