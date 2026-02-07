package orchestration

import (
	"encoding/json"
	"time"
)

type (
	JobStatus string

	JobEventType string

	JobRequest struct {
		ProjectID    string         `json:"project_id"`
		Task         string         `json:"task"`
		Payload      map[string]any `json:"payload"`
		TimeoutSecs  int            `json:"timeout_secs,omitempty"`
		Dependencies []string       `json:"dependencies,omitempty"`
	}

	JobResult struct {
		ID        string          `json:"id"`
		Status    JobStatus       `json:"status"`
		Result    json.RawMessage `json:"result,omitempty"`
		Error     string          `json:"error,omitempty"`
		WorkerID  string          `json:"worker_id,omitempty"`
		CreatedAt time.Time       `json:"created_at"`
		StartedAt *time.Time      `json:"started_at,omitempty"`
		EndedAt   *time.Time      `json:"ended_at,omitempty"`
	}

	WorkerCapabilities struct {
		Projects []string `json:"projects"`
		Hostname string   `json:"hostname"`
	}

	JobEvent struct {
		Type    JobEventType    `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
)

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusAssigned   JobStatus = "assigned"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"

	JobEventHandshake JobEventType = "handshake"
	JobEventJob       JobEventType = "job"
	JobEventCompleted JobEventType = "completed"
	JobEventFailed    JobEventType = "failed"
)
