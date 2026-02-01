package matches

type ListMatchesRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,oneof=10 25 50"`
	Keyword  string `form:"keyword" binding:"omitempty,max=64"`
	Q        string `form:"q" binding:"omitempty,max=255"`
	Issuer   string `form:"issuer" binding:"omitempty,max=255"`
	DateFrom string `form:"date_from" binding:"omitempty"`
	DateTo   string `form:"date_to" binding:"omitempty"`
	NewOnly  string `form:"new_only" binding:"omitempty,oneof=true false"`
	Sort     string `form:"sort" binding:"omitempty,oneof=first_seen_desc last_seen_desc domain_asc"`
}

type ListMatchesResponse struct {
	Items []MatchItem `json:"items"`
	Total int         `json:"total"`
}

type MatchItem struct {
	ID           string  `json:"id"`
	KeywordID    string  `json:"keyword_id"`
	KeywordValue string  `json:"keyword_value"`
	CertSHA256   string  `json:"certificate_sha256"`
	MatchedField string  `json:"matched_field"`
	MatchedValue string  `json:"matched_value"`
	DomainName   *string `json:"domain_name"`
	IssuerCN     *string `json:"issuer_cn"`
	IssuerOrg    *string `json:"issuer_org"`
	SubjectCN    *string `json:"subject_cn"`
	SubjectOrg   *string `json:"subject_org"`
	NotBefore    *string `json:"not_before"`
	NotAfter     *string `json:"not_after"`
	FirstSeenAt  string  `json:"first_seen_at"`
	LastSeenAt   string  `json:"last_seen_at"`
	IsNew        bool    `json:"is_new"`
	CtLogIndex   int64   `json:"ct_log_index"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
