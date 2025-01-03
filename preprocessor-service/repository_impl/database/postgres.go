package database

import (
	"context"
	"database/sql"

	"github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/database/postgres"
	"github.com/anuragprog/notyoutube/preprocessor-service/types/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresDatabase struct {
	db      *sql.DB
	queries *postgres.Queries
}

func NewPostgresDatabase(db *sql.DB, queries *postgres.Queries) *PostgresDatabase {
	return &PostgresDatabase{db, queries}
}

func (pd *PostgresDatabase) Close() error {
	return nil
}

func (pd *PostgresDatabase) CreateDAG(ctx context.Context, dag database.Dag) error {
	dagId := pgtype.UUID{Bytes: dag.ID, Valid: true}
	_, err := pd.queries.CreateDAG(ctx, dagId)
	return err
}

func (pd *PostgresDatabase) GetDAG(ctx context.Context, id uuid.UUID) (database.Dag, error) {
	dagId := pgtype.UUID{Bytes: id, Valid: true}
	dag, err := pd.queries.GetDAG(ctx, dagId)
	if err != nil {
		return database.Dag{}, err
	}
	return database.Dag{
		ID: dag.ID.Bytes,
		CreatedAt: dag.CreatedAt.Time,
	}, nil
}

func (pd *PostgresDatabase) CreateWorkers(ctx context.Context, workers []database.Worker) error {

}

func (pd *PostgresDatabase) ListWorkersOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.Worker, error)
func (pd *PostgresDatabase) CreateDependencies(ctx context.Context, dependencies []database.Dependency) error
func (pd *PostgresDatabase) CreateDependencySources(ctx context.Context, sources []database.DependencySource) error
func (pd *PostgresDatabase) CreateDependencyTargets(ctx context.Context, targets []database.DependencyTarget) error
func (pd *PostgresDatabase) ListDependenciesOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.Dependency, error)
func (pd *PostgresDatabase) ListDependencySourcesOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.DependencySource, error)
func (pd *PostgresDatabase) ListDependencyTargetsOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.DependencySource, error)
func (pd *PostgresDatabase) WithTransaction(ctx context.Context, fn func(repo Database) error) error
