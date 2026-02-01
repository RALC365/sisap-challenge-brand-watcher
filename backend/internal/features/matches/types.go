package matches

import "time"

type SortOrder string

const (
	SortFirstSeenDesc SortOrder = "first_seen_desc"
	SortLastSeenDesc  SortOrder = "last_seen_desc"
	SortDomainAsc     SortOrder = "domain_asc"
)

type Match struct {
	ID              string    `json:"id"`
	KeywordID       string    `json:"keyword_id"`
	KeywordValue    string    `json:"keyword_value"`
	CertSHA256      string    `json:"certificate_sha256"`
	MatchedField    string    `json:"matched_field"`
	MatchedValue    string    `json:"matched_value"`
	DomainName      *string   `json:"domain_name"`
	IssuerCN        *string   `json:"issuer_cn"`
	IssuerOrg       *string   `json:"issuer_org"`
	SubjectCN       *string   `json:"subject_cn"`
	SubjectOrg      *string   `json:"subject_org"`
	NotBefore       *time.Time `json:"not_before"`
	NotAfter        *time.Time `json:"not_after"`
	FirstSeenAt     time.Time `json:"first_seen_at"`
	LastSeenAt      time.Time `json:"last_seen_at"`
	IsNew           bool      `json:"is_new"`
	CtLogIndex      int64     `json:"ct_log_index"`
}

type MatchRow struct {
	ID              string
	KeywordID       string
	KeywordValue    string
	CertSHA256      string
	MatchedField    string
	MatchedValue    string
	DomainName      *string
	IssuerCN        *string
	IssuerOrg       *string
	SubjectCN       *string
	SubjectOrg      *string
	NotBefore       *time.Time
	NotAfter        *time.Time
	FirstSeenAt     time.Time
	LastSeenAt      time.Time
	IsNew           bool
	CtLogIndex      int64
}

type ListQuery struct {
	Page     int
	PageSize int
	Keyword  string
	Q        string
	Issuer   string
	DateFrom *time.Time
	DateTo   *time.Time
	NewOnly  bool
	Sort     SortOrder
}

type ListResponse struct {
	Items []Match `json:"items"`
	Total int     `json:"total"`
}
