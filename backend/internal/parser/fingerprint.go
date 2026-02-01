package parser

import (
	"crypto/sha256"
	"encoding/hex"
)

func ComputeFingerprint(certDER []byte) string {
	hash := sha256.Sum256(certDER)
	return hex.EncodeToString(hash[:])
}
