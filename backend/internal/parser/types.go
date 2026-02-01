package parser

import "time"

type ParsedCertificate struct {
	Fingerprint string
	SubjectCN   string
	SubjectOrg  string
	IssuerCN    string
	IssuerOrg   string
	SANs        []string
	NotBefore   *time.Time
	NotAfter    *time.Time
	RawDER      []byte
}

type ParseResult struct {
	Certificate *ParsedCertificate
	Error       error
	Index       int64
}
