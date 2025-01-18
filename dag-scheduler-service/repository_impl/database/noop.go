package database

import (
	"context"
	"time"

	dbRepo "github.com/anuragprog/notyoutube/dag-scheduler-service/repository/database"
	dbType "github.com/anuragprog/notyoutube/dag-scheduler-service/types/database"
	"github.com/google/uuid"
)

type NoopDatabase struct{}

func NewNoopDatabase() *NoopDatabase {
	return &NoopDatabase{}
}

func (nd *NoopDatabase) Close() error { return nil }

func (nd *NoopDatabase) CreateDAGState(ctx context.Context, dagState dbType.CreateWorkerStateRequest) (dbType.DagState, error) { return dbType.DagState{}, nil }
func (nd *NoopDatabase) UpdateDAGStateStartTime(ctx context.Context, dagStateId uuid.UUID, startTime time.Time) error { return nil }
func (nd *NoopDatabase) UpdateDAGStateEndTime(ctx context.Context, dagStateId uuid.UUID, endTime time.Time) error { return nil }
func (nd *NoopDatabase) UpdateDAGStateFailureReason(ctx context.Context, dagStateId uuid.UUID, failureReason string) error { return nil }
func (nd *NoopDatabase) GetCurrentDAGState(ctx context.Context, dagId uuid.UUID) (dbType.DagState, error) { return dbType.DagState{}, nil }
func (nd *NoopDatabase) ListDAGStates(ctx context.Context, dagId uuid.UUID) ([]dbType.DagState, error) { return nil, nil }
func (nd *NoopDatabase) CreateWorkerStates(ctx context.Context, workerStates []dbType.CreateWorkerStateRequest) error { return nil }
func (nd *NoopDatabase) UpdateWorkerStateStartTime(ctx context.Context, workerStateId uuid.UUID, startTime time.Time) error { return nil }
func (nd *NoopDatabase) UpdateWorkerStateEndTime(ctx context.Context, workerStateId uuid.UUID, endTime time.Time) error { return nil }
func (nd *NoopDatabase) UpdateWorkerStateFailureReason(ctx context.Context, workerStateId uuid.UUID, failureReason string) error { return nil }
func (nd *NoopDatabase) IncrementWorkerStateRetryCount(ctx context.Context, workerStateId uuid.UUID) error { return nil }
func (nd *NoopDatabase) ListWorkerStatesOfWorker(ctx context.Context, workerId uuid.UUID) ([]dbType.WorkerState, error) { return nil, nil }
func (nd *NoopDatabase) GetCurrentWorkerStateOfWorker(ctx context.Context, workerId uuid.UUID) (dbType.WorkerState, error) { return dbType.WorkerState{}, nil }
func (nd *NoopDatabase) WithTransaction(ctx context.Context, fn func(repo dbRepo.Database) error) error { return nil }
