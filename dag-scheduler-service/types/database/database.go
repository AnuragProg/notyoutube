package database

import (
	"time"

	"github.com/google/uuid"
)

type DagStatus string

const (
	DagStatusPending   DagStatus = "pending"
	DagStatusRunning   DagStatus = "running"
	DagStatusCompleted DagStatus = "completed"
	DagStatusFailed    DagStatus = "failed"
)

type WorkerStatus string

const (
	WorkerStatusPending   WorkerStatus = "pending"
	WorkerStatusRunning   WorkerStatus = "running"
	WorkerStatusCompleted WorkerStatus = "completed"
	WorkerStatusFailed    WorkerStatus = "failed"
)

type DagState struct {
	ID            uuid.UUID  `json:"id"`
	DagID         uuid.UUID  `json:"dag_id"`
	DagStatus     DagStatus  `json:"dag_status"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	FailureReason *string    `json:"failure_reason"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CreateDagStateRequest struct {
	ID        uuid.UUID
	DagID     uuid.UUID
	DagStatus DagStatus
}

type WorkerState struct {
	ID            uuid.UUID    `json:"id"`
	DagID         uuid.UUID    `json:"dag_id"`
	WorkerID      uuid.UUID    `json:"worker_id"`
	WorkerStatus  WorkerStatus `json:"worker_status"`
	StartTime     *time.Time   `json:"start_time"`
	EndTime       *time.Time   `json:"end_time"`
	RetryCount    int          `json:"retry_count"`
	FailureReason *string      `json:"failure_reason"`
	CreatedAt     time.Time    `json:"created_at"`
}

type CreateWorkerStateRequest struct {
	ID           uuid.UUID
	DagID        uuid.UUID
	WorkerID     uuid.UUID
	WorkerStatus WorkerStatus
}
