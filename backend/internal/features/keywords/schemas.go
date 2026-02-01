package keywords

type CreateKeywordRequest struct {
	Value string `json:"value" binding:"required,min=1,max=64"`
}

type ListKeywordsResponse struct {
	Items []KeywordItem `json:"items"`
	Total int           `json:"total"`
}

type KeywordItem struct {
	KeywordID       string `json:"keyword_id"`
	Value           string `json:"value"`
	NormalizedValue string `json:"normalized_value"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
}

type DeleteKeywordResponse struct {
	OK bool `json:"ok"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
