// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package postgres

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DagStatus string

const (
	DagStatusPending   DagStatus = "pending"
	DagStatusRunning   DagStatus = "running"
	DagStatusCompleted DagStatus = "completed"
	DagStatusFailed    DagStatus = "failed"
)

func (e *DagStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = DagStatus(s)
	case string:
		*e = DagStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for DagStatus: %T", src)
	}
	return nil
}

type NullDagStatus struct {
	DagStatus DagStatus `json:"dag_status"`
	Valid     bool      `json:"valid"` // Valid is true if DagStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullDagStatus) Scan(value interface{}) error {
	if value == nil {
		ns.DagStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.DagStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullDagStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.DagStatus), nil
}

type WorkerStatus string

const (
	WorkerStatusPending   WorkerStatus = "pending"
	WorkerStatusRunning   WorkerStatus = "running"
	WorkerStatusCompleted WorkerStatus = "completed"
	WorkerStatusFailed    WorkerStatus = "failed"
)

func (e *WorkerStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = WorkerStatus(s)
	case string:
		*e = WorkerStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for WorkerStatus: %T", src)
	}
	return nil
}

type NullWorkerStatus struct {
	WorkerStatus WorkerStatus `json:"worker_status"`
	Valid        bool         `json:"valid"` // Valid is true if WorkerStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullWorkerStatus) Scan(value interface{}) error {
	if value == nil {
		ns.WorkerStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.WorkerStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullWorkerStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.WorkerStatus), nil
}

type DagState struct {
	ID            uuid.UUID        `json:"id"`
	DagID         uuid.UUID        `json:"dag_id"`
	DagStatus     DagStatus        `json:"dag_status"`
	StartTime     pgtype.Timestamp `json:"start_time"`
	EndTime       pgtype.Timestamp `json:"end_time"`
	FailureReason pgtype.Text      `json:"failure_reason"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
}

type WorkerState struct {
	ID            uuid.UUID        `json:"id"`
	DagID         uuid.UUID        `json:"dag_id"`
	WorkerID      uuid.UUID        `json:"worker_id"`
	WorkerStatus  WorkerStatus     `json:"worker_status"`
	StartTime     pgtype.Timestamp `json:"start_time"`
	EndTime       pgtype.Timestamp `json:"end_time"`
	RetryCount    pgtype.Int4      `json:"retry_count"`
	FailureReason pgtype.Text      `json:"failure_reason"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
}
