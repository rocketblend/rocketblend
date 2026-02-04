package types

import (
	"encoding/json"
	"time"
)

type (
	// JobStatus represents the current state of a job
	JobStatus string

	// JobRequest is the payload sent by clients to submit work
	JobRequest struct {
		ProjectID    string                 `json:"project_id"`
		Task         string                 `json:"task"`
		Payload      map[string]interface{} `json:"payload"`
		TimeoutSecs  int                    `json:"timeout_secs,omitempty"`
		Dependencies []string               `json:"dependencies,omitempty"`
	}

	// JobResult represents the outcome of a job execution
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

	// WorkerCapabilities describes what a worker can do
	WorkerCapabilities struct {
		Projects []string `json:"projects"` // Project IDs (folder names) this worker manages
		Hostname string   `json:"hostname"`
	}

	// JobEvent is pushed via WebSocket or Polling
	JobEvent struct {
		Type    string          `json:"type"` // "assigned", "completed", "failed"
		Payload json.RawMessage `json:"payload"`
	}
)

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusAssigned   JobStatus = "assigned"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)
