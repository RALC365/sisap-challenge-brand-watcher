package parser

import (
	"crypto/x509"
	"encoding/base64"
	"errors"

	"brand-protection-monitor/internal/ct"
)

var (
	ErrInvalidLeafInput = errors.New("invalid leaf input")
	ErrNoCertificate    = errors.New("no certificate in entry")
	ErrParseFailed      = errors.New("certificate parse failed")
)

func ParseEntry(entry ct.LogEntry, index int64) ParseResult {
	leafBytes, err := base64.StdEncoding.DecodeString(entry.LeafInput)
	if err != nil {
		return ParseResult{Error: ErrInvalidLeafInput, Index: index}
	}

	if len(leafBytes) < 15 {
		return ParseResult{Error: ErrInvalidLeafInput, Index: index}
	}

	entryType := uint16(leafBytes[10])<<8 | uint16(leafBytes[11])
	if entryType == 1 {
		return parsePrecert(leafBytes, index)
	}

	return parseX509Entry(leafBytes, index)
}

func parseX509Entry(leafBytes []byte, index int64) ParseResult {
	if len(leafBytes) < 15 {
		return ParseResult{Error: ErrInvalidLeafInput, Index: index}
	}

	certStart := 15
	if certStart >= len(leafBytes) {
		return ParseResult{Error: ErrNoCertificate, Index: index}
	}

	certLen := int(leafBytes[12])<<16 | int(leafBytes[13])<<8 | int(leafBytes[14])
	if certStart+certLen > len(leafBytes) {
		return ParseResult{Error: ErrInvalidLeafInput, Index: index}
	}

	certDER := leafBytes[certStart : certStart+certLen]
	return parseCertDER(certDER, index)
}

func parsePrecert(leafBytes []byte, index int64) ParseResult {
	if len(leafBytes) < 47 {
		return ParseResult{Error: ErrInvalidLeafInput, Index: index}
	}

	tbsStart := 47
	if tbsStart >= len(leafBytes) {
		return ParseResult{Error: ErrNoCertificate, Index: index}
	}

	tbsLen := int(leafBytes[44])<<16 | int(leafBytes[45])<<8 | int(leafBytes[46])
	if tbsStart+tbsLen > len(leafBytes) {
		return ParseResult{Error: ErrInvalidLeafInput, Index: index}
	}

	tbsDER := leafBytes[tbsStart : tbsStart+tbsLen]

	cert, err := x509.ParseCertificate(tbsDER)
	if err != nil {
		return parseCertFromExtraData(leafBytes, index)
	}

	return buildParsedCert(cert, tbsDER, index)
}

func parseCertFromExtraData(leafBytes []byte, index int64) ParseResult {
	return ParseResult{Error: ErrParseFailed, Index: index}
}

func parseCertDER(certDER []byte, index int64) ParseResult {
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return ParseResult{Error: ErrParseFailed, Index: index}
	}

	return buildParsedCert(cert, certDER, index)
}

func buildParsedCert(cert *x509.Certificate, rawDER []byte, index int64) ParseResult {
	parsed := &ParsedCertificate{
		Fingerprint: ComputeFingerprint(rawDER),
		SubjectCN:   NormalizeDomain(cert.Subject.CommonName),
		IssuerCN:    cert.Issuer.CommonName,
		RawDER:      rawDER,
	}

	if len(cert.Subject.Organization) > 0 {
		parsed.SubjectOrg = cert.Subject.Organization[0]
	}
	if len(cert.Issuer.Organization) > 0 {
		parsed.IssuerOrg = cert.Issuer.Organization[0]
	}

	parsed.SANs = NormalizeDomains(cert.DNSNames)

	notBefore := cert.NotBefore
	notAfter := cert.NotAfter
	parsed.NotBefore = &notBefore
	parsed.NotAfter = &notAfter

	return ParseResult{Certificate: parsed, Index: index}
}
