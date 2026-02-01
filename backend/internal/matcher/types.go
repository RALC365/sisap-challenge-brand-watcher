package matcher

type MatchedField string

const (
	MatchedFieldCN   MatchedField = "cn"
	MatchedFieldSAN  MatchedField = "san"
	MatchedFieldBoth MatchedField = "both"
)

type Keyword struct {
	ID              string
	Value           string
	NormalizedValue string
}

type Match struct {
	KeywordID    string
	KeywordValue string
	MatchedField MatchedField
	MatchedValue string
	DomainName   string
}
