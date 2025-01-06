package database

import (
	"context"
	"fmt"
	"time"

	dbRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/database"
	"github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/database/postgres"
	dbType "github.com/anuragprog/notyoutube/preprocessor-service/types/database"
	"github.com/anuragprog/notyoutube/preprocessor-service/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var WorkerTypeToPostgresWorkerType = map[dbType.WorkerType]postgres.WorkerType{
	dbType.WorkerTypeVideoEncoder:       postgres.WorkerTypeVideoEncoder,
	dbType.WorkerTypeAsciiEncoder:       postgres.WorkerTypeAsciiEncoder,
	dbType.WorkerTypeThumbnailGenerator: postgres.WorkerTypeThumbnailGenerator,
	dbType.WorkerTypeAssembler:          postgres.WorkerTypeAssembler,
	dbType.WorkerTypeVideoExtractor:     postgres.WorkerTypeVideoExtractor,
	dbType.WorkerTypeAudioExtractor:     postgres.WorkerTypeAudioExtractor,
	dbType.WorkerTypeMetadataExtractor:  postgres.WorkerTypeMetadataExtractor,
}
var postgresWorkerTypeToWorkerType = map[postgres.WorkerType]dbType.WorkerType{
	postgres.WorkerTypeVideoEncoder:       dbType.WorkerTypeVideoEncoder,
	postgres.WorkerTypeAsciiEncoder:       dbType.WorkerTypeAsciiEncoder,
	postgres.WorkerTypeThumbnailGenerator: dbType.WorkerTypeThumbnailGenerator,
	postgres.WorkerTypeAssembler:          dbType.WorkerTypeAssembler,
	postgres.WorkerTypeVideoExtractor:     dbType.WorkerTypeVideoExtractor,
	postgres.WorkerTypeAudioExtractor:     dbType.WorkerTypeAudioExtractor,
	postgres.WorkerTypeMetadataExtractor:  dbType.WorkerTypeMetadataExtractor,
}

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

func (pd *PostgresDatabase) CreateDAG(ctx context.Context, dag dbType.Dag) error {
	dagId := pgtype.UUID{Bytes: dag.ID, Valid: true}
	_, err := pd.queries.CreateDAG(ctx, dagId)
	return err
}

func (pd *PostgresDatabase) GetDAG(ctx context.Context, id uuid.UUID) (dbType.Dag, error) {
	dagId := pgtype.UUID{Bytes: id, Valid: true}
	dag, err := pd.queries.GetDAG(ctx, dagId)
	if err != nil {
		return dbType.Dag{}, err
	}
	return dbType.Dag{
		ID:        dag.ID.Bytes,
		CreatedAt: dag.CreatedAt.Time,
	}, nil
}

func (pd *PostgresDatabase) CreateWorkers(ctx context.Context, workers []dbType.Worker) error {
	result := pd.queries.CreateWorkers(
		ctx,
		utils.Map(workers, func(w dbType.Worker) postgres.CreateWorkersParams {
			return postgres.CreateWorkersParams{
				ID:           pgtype.UUID{Bytes: w.ID, Valid: true},
				DagID:        pgtype.UUID{Bytes: w.DagID, Valid: true},
				Name:         w.Name,
				Description:  pgtype.Text{String: w.Description, Valid: true},
				WorkerType:   WorkerTypeToPostgresWorkerType[w.WorkerType],
				WorkerConfig: w.WorkerConfig,
			}
		}),
	)
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

func (pd *PostgresDatabase) ListWorkersOfDAG(ctx context.Context, dagId uuid.UUID) ([]dbType.Worker, error) {
	workers, err := pd.queries.ListWorkersOfDAG(ctx, pgtype.UUID{Bytes: dagId, Valid: true})
	if err != nil {
		return nil, err
	}

	return utils.Map(workers, func(w postgres.Worker) dbType.Worker {
		return dbType.Worker{
			ID:           w.ID.Bytes,
			DagID:        w.DagID.Bytes,
			Name:         w.Name,
			Description:  w.Description.String,
			WorkerType:   postgresWorkerTypeToWorkerType[w.WorkerType],
			WorkerConfig: w.WorkerConfig,
		}
	}), nil
}

func (pd *PostgresDatabase) CreateDependencies(ctx context.Context, dependencies []dbType.Dependency) error {
	result := pd.queries.CreateDependencies(
		ctx,
		utils.Map(
			dependencies,
			func(d dbType.Dependency) postgres.CreateDependenciesParams {
				return postgres.CreateDependenciesParams{
					ID:    pgtype.UUID{Bytes: d.ID, Valid: true},
					DagID: pgtype.UUID{Bytes: d.DagID, Valid: true},
				}
			},
		),
	)
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

func (pd *PostgresDatabase) CreateDependencySources(ctx context.Context, sources []dbType.DependencySource) error {
	result := pd.queries.CreateDependencySources(
		ctx,
		utils.Map(
			sources,
			func(ds dbType.DependencySource) postgres.CreateDependencySourcesParams {
				return postgres.CreateDependencySourcesParams{
					ID:           pgtype.UUID{Bytes: ds.ID, Valid: true},
					DagID:        pgtype.UUID{Bytes: ds.DagID, Valid: true},
					DependencyID: pgtype.UUID{Bytes: ds.DependencyID, Valid: true},
					SourceID:     pgtype.UUID{Bytes: ds.SourceID, Valid: true},
				}
			},
		),
	)
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

func (pd *PostgresDatabase) CreateDependencyTargets(ctx context.Context, targets []dbType.DependencyTarget) error {
	result := pd.queries.CreateDependencyTargets(
		ctx,
		utils.Map(
			targets,
			func(dt dbType.DependencyTarget) postgres.CreateDependencyTargetsParams {
				return postgres.CreateDependencyTargetsParams{
					ID:           pgtype.UUID{Bytes: dt.ID, Valid: true},
					DagID:        pgtype.UUID{Bytes: dt.DagID, Valid: true},
					DependencyID: pgtype.UUID{Bytes: dt.DependencyID, Valid: true},
					TargetID:     pgtype.UUID{Bytes: dt.TargetID, Valid: true},
				}
			},
		),
	)

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

func (pd *PostgresDatabase) ListDependenciesOfDAG(ctx context.Context, dagId uuid.UUID) ([]dbType.Dependency, error) {
	deps, err := pd.queries.ListDependenciesOfDAG(
		ctx,
		pgtype.UUID{Bytes: dagId, Valid: true},
	)
	if err != nil {
		return nil, err
	}

	return utils.Map(
		deps,
		func(d postgres.Dependency) dbType.Dependency {
			return dbType.Dependency{
				ID:    d.ID.Bytes,
				DagID: d.DagID.Bytes,
			}
		},
	), nil
}

func (pd *PostgresDatabase) ListDependencySourcesOfDAG(ctx context.Context, dagId uuid.UUID) ([]dbType.DependencySource, error) {
	depSrcs, err := pd.queries.ListDependencySourcesOfDependency(
		ctx,
		pgtype.UUID{Bytes: dagId, Valid: true},
	)
	if err != nil {
		return nil, err
	}
	return utils.Map(
		depSrcs,
		func(ds postgres.DependencySource) dbType.DependencySource {
			return dbType.DependencySource{
				ID:           ds.ID.Bytes,
				DagID:        ds.DagID.Bytes,
				DependencyID: ds.DependencyID.Bytes,
				SourceID:     ds.SourceID.Bytes,
			}
		},
	), nil
}

func (pd *PostgresDatabase) ListDependencyTargetsOfDAG(ctx context.Context, dagId uuid.UUID) ([]dbType.DependencyTarget, error) {
	depTgts, err := pd.queries.ListDependencyTargetsOfDependency(
		ctx,
		pgtype.UUID{Bytes: dagId, Valid: true},
	)
	if err != nil {
		return nil, err
	}
	return utils.Map(
		depTgts,
		func(dt postgres.DependencyTarget) dbType.DependencyTarget {
			return dbType.DependencyTarget{
				ID:           dt.ID.Bytes,
				DagID:        dt.DagID.Bytes,
				DependencyID: dt.DependencyID.Bytes,
				TargetID:     dt.TargetID.Bytes,
			}
		},
	), nil
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
