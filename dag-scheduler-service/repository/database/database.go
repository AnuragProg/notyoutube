package database

import (
	"context"
	"io"
	"time"

	"github.com/anuragprog/notyoutube/dag-scheduler-service/types/database"
	"github.com/google/uuid"
)

type Database interface {

	// to make sure the handler closes the connection properly
	io.Closer

	CreateDAGState(ctx context.Context, dagState database.CreateWorkerStateRequest) (database.DagState, error)
	UpdateDAGStateStartTime(ctx context.Context, dagStateId uuid.UUID, startTime time.Time) error
	UpdateDAGStateEndTime(ctx context.Context, dagStateId uuid.UUID, endTime time.Time) error
	UpdateDAGStateFailureReason(ctx context.Context, dagStateId uuid.UUID, failureReason string) error
	GetCurrentDAGState(ctx context.Context, dagId uuid.UUID) (database.DagState, error)
	ListDAGStates(ctx context.Context, dagId uuid.UUID) ([]database.DagState, error)

	CreateWorkerStates(ctx context.Context, workerStates []database.CreateWorkerStateRequest) error 
	UpdateWorkerStateStartTime(ctx context.Context, workerStateId uuid.UUID, startTime time.Time) error
	UpdateWorkerStateEndTime(ctx context.Context, workerStateId uuid.UUID, endTime time.Time) error
	UpdateWorkerStateFailureReason(ctx context.Context, workerStateId uuid.UUID, failureReason string) error
	ListWorkerStatesOfWorker(ctx context.Context, workerId uuid.UUID) ([]database.WorkerState, error)
	GetCurrentWorkerStateOfWorker(ctx context.Context, workerId uuid.UUID) (database.WorkerState, error)

	WithTransaction(ctx context.Context, fn func(repo Database) error) error
}
