package ct

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	config     ClientConfig
	httpClient *http.Client
}

func NewClient(config ClientConfig) *Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: config.ConnectTimeout,
		}).DialContext,
		ResponseHeaderTimeout: config.ReadTimeout,
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   config.ConnectTimeout + config.ReadTimeout,
		},
	}
}

func (c *Client) GetSTH(ctx context.Context) (*STHResponse, error) {
	url := fmt.Sprintf("%s/ct/v1/get-sth", c.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, ErrTimeout
		}
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrHTTPError, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read body: %v", ErrConnectionError, err)
	}

	var sth STHResponse
	if err := json.Unmarshal(body, &sth); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return &sth, nil
}

func (c *Client) GetEntries(ctx context.Context, start, end int64) (*GetEntriesResponse, error) {
	url := fmt.Sprintf("%s/ct/v1/get-entries?start=%d&end=%d", c.config.BaseURL, start, end)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, ErrTimeout
		}
		return nil, fmt.Errorf("%w: %v", ErrConnectionError, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrHTTPError, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read body: %v", ErrConnectionError, err)
	}

	var entries GetEntriesResponse
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return &entries, nil
}

func (c *Client) GetEntriesChunked(ctx context.Context, start, end int64, chunkSize int) ([]LogEntry, error) {
	var allEntries []LogEntry

	for current := start; current <= end; current += int64(chunkSize) {
		chunkEnd := current + int64(chunkSize) - 1
		if chunkEnd > end {
			chunkEnd = end
		}

		resp, err := c.GetEntries(ctx, current, chunkEnd)
		if err != nil {
			return allEntries, err
		}

		allEntries = append(allEntries, resp.Entries...)

		if chunkEnd >= end {
			break
		}
	}

	return allEntries, nil
}

func CalculateRange(treeSize int64, batchSize int, lastProcessedIndex int64) FetchRange {
	if treeSize == 0 {
		return FetchRange{Start: 0, End: 0}
	}

	end := treeSize - 1

	var start int64
	if lastProcessedIndex > 0 {
		start = lastProcessedIndex + 1
	} else {
		start = end - int64(batchSize) + 1
		if start < 0 {
			start = 0
		}
	}

	if start > end {
		return FetchRange{Start: end, End: end}
	}

	if end-start+1 > int64(batchSize) {
		end = start + int64(batchSize) - 1
	}

	return FetchRange{Start: start, End: end}
}

var _ = time.Second
