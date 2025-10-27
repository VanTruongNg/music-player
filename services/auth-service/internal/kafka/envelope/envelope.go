package envelope

import (
	"encoding/json"
	"time"
)

type Envelope struct {
	MessageID string          `json:"message_id"`
	Version   string          `json:"version"`
	Timestamp time.Time       `json:"timestamp"`
	Source    string          `json:"source"`
	Priority  Priority        `json:"priority"`
	Data      json.RawMessage `json:"data"`
	Metadata  *Metadata       `json:"metadata,omitempty"`
	TraceID   string          `json:"trace_id,omitempty"`
	SpanID    string          `json:"span_id,omitempty"`
}

type Metadata struct {
	UserID       string            `json:"user_id,omitempty"`
	TenantID     string            `json:"tenant_id,omitempty"`
	Environment  string            `json:"environment,omitempty"`
	RetryCount   int               `json:"retry_count,omitempty"`
	CustomFields map[string]string `json:"custom_fields,omitempty"`
}

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityNormal   Priority = "normal"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityNormal, PriorityHigh, PriorityCritical:
		return true
	default:
		return false
	}
}

func (p Priority) String() string {
	return string(p)
}
