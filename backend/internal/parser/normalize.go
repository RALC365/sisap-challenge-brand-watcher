package parser

import "strings"

func NormalizeDomain(domain string) string {
	domain = strings.ToLower(domain)
	domain = strings.TrimSuffix(domain, ".")
	return domain
}

func NormalizeDomains(domains []string) []string {
	normalized := make([]string, 0, len(domains))
	seen := make(map[string]bool)

	for _, d := range domains {
		n := NormalizeDomain(d)
		if n != "" && !seen[n] {
			normalized = append(normalized, n)
			seen[n] = true
		}
	}

	return normalized
}
