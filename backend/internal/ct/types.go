package ct

import (
	"errors"
	"time"
)

var (
	ErrInvalidJSON     = errors.New("invalid JSON response from CT log")
	ErrHTTPError       = errors.New("HTTP error from CT log")
	ErrConnectionError = errors.New("connection error to CT log")
	ErrTimeout         = errors.New("timeout connecting to CT log")
)

type STHResponse struct {
	TreeSize          int64  `json:"tree_size"`
	Timestamp         int64  `json:"timestamp"`
	SHA256RootHash    string `json:"sha256_root_hash"`
	TreeHeadSignature string `json:"tree_head_signature"`
}

type GetEntriesResponse struct {
	Entries []LogEntry `json:"entries"`
}

type LogEntry struct {
	LeafInput string `json:"leaf_input"`
	ExtraData string `json:"extra_data"`
}

type ClientConfig struct {
	BaseURL        string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	BatchSize      int
}

type FetchRange struct {
	Start int64
	End   int64
}
