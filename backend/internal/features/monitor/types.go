package monitor

import "time"

type MonitorState string

const (
        MonitorStateIdle    MonitorState = "idle"
        MonitorStateRunning MonitorState = "running"
        MonitorStateError   MonitorState = "error"
)

type MetricsLastRun struct {
        ProcessedCount  int `json:"processed_count"`
        MatchCount      int `json:"match_count"`
        ParseErrorCount int `json:"parse_error_count"`
        DurationMs      int `json:"duration_ms"`
        CtLatencyMs     int `json:"ct_latency_ms"`
        DbLatencyMs     int `json:"db_latency_ms"`
}

type StatusResponse struct {
        State               MonitorState    `json:"state"`
        LastRunAt           *time.Time      `json:"last_run_at"`
        LastSuccessAt       *time.Time      `json:"last_success_at"`
        LastErrorCode       *string         `json:"last_error_code"`
        LastErrorMessage    *string         `json:"last_error_message"`
        MetricsLastRun      *MetricsLastRun `json:"metrics_last_run"`
        PollIntervalSeconds int             `json:"poll_interval_seconds"`
}

type MonitorStateRow struct {
        ID               int
        State            string
        LastRunAt        *time.Time
        LastSuccessAt    *time.Time
        LastErrorCode    *string
        LastErrorMessage *string
}

type LastRunRow struct {
        CertificatesProcessed int
        MatchesFound          int
        ParseErrorCount       int
        DurationMs            *int
        CtLatencyMs           *int
        DbLatencyMs           *int
}
