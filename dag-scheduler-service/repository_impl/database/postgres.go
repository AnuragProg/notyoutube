package database

import (
	"context"
	"fmt"
	"time"

	"github.com/anuragprog/notyoutube/dag-scheduler-service/utils"
	dbType "github.com/anuragprog/notyoutube/dag-scheduler-service/types/database"
	dbRepo "github.com/anuragprog/notyoutube/dag-scheduler-service/repository/database"
	"github.com/anuragprog/notyoutube/dag-scheduler-service/repository_impl/database/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDatabase struct {
	pool    *pgxpool.Pool
	queries *postgres.Queries
}

func NewPostgresDatabase(ctx context.Context, host, port, user, password, dbname string) (*PostgresDatabase, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)
	// db *sql.DB, queries *postgres.Queries
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	queries := postgres.New(pool)
	return &PostgresDatabase{pool, queries}, nil
}

func (pd *PostgresDatabase) Close() error {
	pd.pool.Close()
	return nil
}

func (pd *PostgresDatabase) CreateDAGState(ctx context.Context, dagState dbType.CreateWorkerStateRequest) (dbType.DagState, error) {
	createdDAGState, err := pd.queries.CreateDAGState(ctx, postgres.CreateDAGStateParams{
		ID: dagState.ID,
		DagID: dagState.DagID,
		DagStatus: postgres.DagStatus(dagState.WorkerStatus),
	})
	if err != nil {
		return dbType.DagState{}, nil
	}

	var startTime, endTime *time.Time
	var failureReason *string

	if createdDAGState.StartTime.Valid {
		startTime = &createdDAGState.StartTime.Time
	}
	if createdDAGState.EndTime.Valid {
		endTime = &createdDAGState.EndTime.Time
	}
	if createdDAGState.FailureReason.Valid {
		failureReason = &createdDAGState.FailureReason.String
	}

	return dbType.DagState{
		ID: createdDAGState.ID,
		DagID: createdDAGState.DagID,
		DagStatus: dbType.DagStatus(createdDAGState.DagStatus),
		StartTime: startTime,
		EndTime: endTime,
		FailureReason: failureReason,
		CreatedAt: createdDAGState.CreatedAt.Time,
	}, nil
}

func (pd *PostgresDatabase) UpdateDAGStateStartTime(ctx context.Context, dagStateId uuid.UUID, startTime time.Time) error {
	return pd.queries.UpdateDAGStateStartTime(ctx, postgres.UpdateDAGStateStartTimeParams{
		ID: dagStateId,
		StartTime: pgtype.Timestamp{
			Time: startTime,
			Valid: true,
		},
	})
}
	
func (pd *PostgresDatabase) UpdateDAGStateEndTime(ctx context.Context, dagStateId uuid.UUID, endTime time.Time) error {
	return pd.queries.UpdateDAGStateEndTime(ctx, postgres.UpdateDAGStateEndTimeParams{
		ID: dagStateId,
		EndTime: pgtype.Timestamp{
			Time: endTime,
			Valid: true,
		},
	})
}

func (pd *PostgresDatabase) UpdateDAGStateFailureReason(ctx context.Context, dagStateId uuid.UUID, failureReason string) error {
	return pd.queries.UpdateDAGStateFailureReason(ctx, postgres.UpdateDAGStateFailureReasonParams{
		ID: dagStateId,
		FailureReason: pgtype.Text{
			String: failureReason,
			Valid: true,
		},
	})
}

func (pd *PostgresDatabase) GetCurrentDAGState(ctx context.Context, dagId uuid.UUID) (dbType.DagState, error) {
	dagState, err := pd.queries.GetCurrentDAGState(ctx, dagId)
	if err != nil {
		return dbType.DagState{}, nil
	}
	var startTime, endTime *time.Time
	var failureReason *string

	if dagState.StartTime.Valid {
		startTime = &dagState.StartTime.Time
	}
	if dagState.EndTime.Valid {
		endTime = &dagState.EndTime.Time
	}
	if dagState.FailureReason.Valid {
		failureReason = &dagState.FailureReason.String
	}
	return dbType.DagState{
		ID: dagState.ID,
		DagID: dagState.DagID,
		DagStatus: dbType.DagStatus(dagState.DagStatus),
		StartTime: startTime,
		EndTime: endTime,
		FailureReason: failureReason,
		CreatedAt: dagState.CreatedAt.Time,
	}, nil
}

func (pd *PostgresDatabase) ListDAGStates(ctx context.Context, dagId uuid.UUID) ([]dbType.DagState, error) {
	dagStates, err := pd.queries.ListDAGStates(ctx, dagId)
	if err != nil {
		return nil, err
	}

	return utils.Map(dagStates, func(dagState postgres.DagState) dbType.DagState {
		var startTime, endTime *time.Time
		var failureReason *string

		if dagState.StartTime.Valid {
			startTime = &dagState.StartTime.Time
		}
		if dagState.EndTime.Valid {
			endTime = &dagState.EndTime.Time
		}
		if dagState.FailureReason.Valid {
			failureReason = &dagState.FailureReason.String
		}
		return dbType.DagState{
			ID: dagState.ID,
			DagID: dagState.DagID,
			DagStatus: dbType.DagStatus(dagState.DagStatus),
			StartTime: startTime,
			EndTime: endTime,
			FailureReason: failureReason,
			CreatedAt: dagState.CreatedAt.Time,
		}
	}), nil
}

func (pd *PostgresDatabase) CreateWorkerStates(ctx context.Context, workerStates []dbType.CreateWorkerStateRequest) error {
	result := pd.queries.CreateWorkerStates(ctx, utils.Map(
		workerStates,
		func(workerState dbType.CreateWorkerStateRequest) postgres.CreateWorkerStatesParams {
			return postgres.CreateWorkerStatesParams{
				ID: workerState.ID,
				DagID: workerState.DagID,
				WorkerID: workerState.WorkerID,
				WorkerStatus: postgres.WorkerStatus(workerState.WorkerStatus),
			}
		},
	))
	defer result.Close()

	var err error
	result.Exec(func(i int, _err error) {
		if _err != nil {
			err = _err
			result.Close()
			return
		}
	})
	return err
}

func (pd *PostgresDatabase) UpdateWorkerStateStartTime(ctx context.Context, workerStateId uuid.UUID, startTime time.Time) error {
	return pd.queries.UpdateWorkerStateStartTime(ctx, postgres.UpdateWorkerStateStartTimeParams{
		ID: workerStateId,
		StartTime: pgtype.Timestamp{
			Time: startTime,
			Valid: true,
		},
	})
}

func (pd *PostgresDatabase) UpdateWorkerStateEndTime(ctx context.Context, workerStateId uuid.UUID, endTime time.Time) error {
	return pd.queries.UpdateWorkerStateEndTime(ctx, postgres.UpdateWorkerStateEndTimeParams{
		ID: workerStateId,
		EndTime: pgtype.Timestamp{
			Time: endTime,
			Valid: true,
		},
	})
}
func (pd *PostgresDatabase) UpdateWorkerStateFailureReason(ctx context.Context, workerStateId uuid.UUID, failureReason string) error {
	return pd.queries.UpdateWorkerStateFailureReason(ctx, postgres.UpdateWorkerStateFailureReasonParams{
		ID: workerStateId,
		FailureReason: pgtype.Text{
			String: failureReason,
			Valid: true,
		},
	})
}

func (pd *PostgresDatabase) ListWorkerStatesOfWorker(ctx context.Context, workerId uuid.UUID) ([]dbType.WorkerState, error) {
	workerStates, err := pd.queries.ListWorkerStatesOfWorker(ctx, workerId)
	if err != nil {
		return nil, err
	}
	return utils.Map(workerStates, func(workerState postgres.WorkerState) dbType.WorkerState {
		var startTime, endTime *time.Time
		var failureReason *string

		if workerState.StartTime.Valid {
			startTime = &workerState.StartTime.Time
		}
		if workerState.EndTime.Valid {
			endTime = &workerState.EndTime.Time
		}
		if workerState.FailureReason.Valid {
			failureReason = &workerState.FailureReason.String
		}
		return dbType.WorkerState{
			ID: workerState.ID,
			DagID: workerState.DagID,
			WorkerID: workerState.WorkerID,
			WorkerStatus: dbType.WorkerStatus(workerState.WorkerStatus),
			StartTime: startTime,
			EndTime: endTime,
			RetryCount: int(workerState.RetryCount.Int32),
			FailureReason: failureReason,
			CreatedAt: workerState.CreatedAt.Time,
		}
	}), nil
}

func (pd *PostgresDatabase) GetCurrentWorkerStateOfWorker(ctx context.Context, workerId uuid.UUID) (dbType.WorkerState, error) {
	workerState, err := pd.queries.GetCurrentWorkerStateOfWorker(ctx, workerId)
	if err != nil {
		return dbType.WorkerState{}, err
	}

	var startTime, endTime *time.Time
	var failureReason *string

	if workerState.StartTime.Valid {
		startTime = &workerState.StartTime.Time
	}
	if workerState.EndTime.Valid {
		endTime = &workerState.EndTime.Time
	}
	if workerState.FailureReason.Valid {
		failureReason = &workerState.FailureReason.String
	}
	return dbType.WorkerState{
			ID: workerState.ID,
			DagID: workerState.DagID,
			WorkerID: workerState.WorkerID,
			WorkerStatus: dbType.WorkerStatus(workerState.WorkerStatus),
			StartTime: startTime,
			EndTime: endTime,
			RetryCount: int(workerState.RetryCount.Int32),
			FailureReason: failureReason,
			CreatedAt: workerState.CreatedAt.Time,
	}, nil
}

func (pd *PostgresDatabase) WithTransaction(ctx context.Context, fn func(repo dbRepo.Database) error) error {
	tx, err := pd.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// transaction repo
	queries := pd.queries.WithTx(tx)
	var txRepo dbRepo.Database = &PostgresDatabase{pool: pd.pool, queries: queries}
	if err = fn(txRepo); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
