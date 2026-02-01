package matcher

import (
	"strings"

	"brand-protection-monitor/internal/parser"
)

type Matcher struct {
	keywords []Keyword
}

func New(keywords []Keyword) *Matcher {
	return &Matcher{keywords: keywords}
}

func (m *Matcher) Match(cert *parser.ParsedCertificate) []Match {
	var matches []Match

	for _, kw := range m.keywords {
		normalizedKeyword := strings.ToLower(kw.NormalizedValue)
		cnMatch := false
		sanMatch := false
		var matchedSAN string

		if cert.SubjectCN != "" && strings.Contains(strings.ToLower(cert.SubjectCN), normalizedKeyword) {
			cnMatch = true
		}

		for _, san := range cert.SANs {
			if strings.Contains(strings.ToLower(san), normalizedKeyword) {
				sanMatch = true
				matchedSAN = san
				break
			}
		}

		if cnMatch || sanMatch {
			match := Match{
				KeywordID:    kw.ID,
				KeywordValue: kw.Value,
			}

			if cnMatch && sanMatch {
				match.MatchedField = MatchedFieldBoth
				match.MatchedValue = cert.SubjectCN
				match.DomainName = cert.SubjectCN
			} else if cnMatch {
				match.MatchedField = MatchedFieldCN
				match.MatchedValue = cert.SubjectCN
				match.DomainName = cert.SubjectCN
			} else {
				match.MatchedField = MatchedFieldSAN
				match.MatchedValue = matchedSAN
				match.DomainName = matchedSAN
			}

			matches = append(matches, match)
		}
	}

	return matches
}
