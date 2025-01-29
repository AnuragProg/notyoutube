package database

import (
	"context"

	"github.com/anuragprog/notyoutube/preprocessor-service/types/database"
	dbRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/database"
	"github.com/google/uuid"
)

type NoopDatabase struct{}

func NewNoopDatabase() *NoopDatabase {
	return &NoopDatabase{}
}

func (nd *NoopDatabase) Close() error { return nil }

func (nd *NoopDatabase) CreateDAG(ctx context.Context, dag database.Dag) error { return nil }
func (nd *NoopDatabase) GetDAG(ctx context.Context, id uuid.UUID) (database.Dag, error) { return database.Dag{}, nil }

func (nd *NoopDatabase) CreateWorkers(ctx context.Context, workers []database.Worker) error { return nil }
func (nd *NoopDatabase) ListWorkersOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.Worker, error) { return []database.Worker{}, nil }
func (nd *NoopDatabase) GetWorkerById(ctx context.Context, workerId uuid.UUID) (database.Worker, error) { return database.Worker{}, nil }

func (nd *NoopDatabase) CreateDependencies(ctx context.Context, dependencies []database.Dependency) error { return nil }
func (nd *NoopDatabase) CreateDependencySources(ctx context.Context, sources []database.DependencySource) error { return nil }
func (nd *NoopDatabase) CreateDependencyTargets(ctx context.Context, targets []database.DependencyTarget) error { return nil }
func (nd *NoopDatabase) ListDependenciesOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.Dependency, error) { return []database.Dependency{}, nil }
func (nd *NoopDatabase) ListDependencySourcesOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.DependencySource, error) { return []database.DependencySource{}, nil }
func (nd *NoopDatabase) ListDependencyTargetsOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.DependencyTarget, error) { return []database.DependencyTarget{}, nil }
func (nd *NoopDatabase) ListDependencySourcesWhereWorkerIsSource(ctx context.Context, workerId uuid.UUID) ([]database.DependencySource, error) { return []database.DependencySource{}, nil } 
func (nd *NoopDatabase) BatchListDependencySourcesOfDependency(ctx context.Context, dependencyIds []uuid.UUID) ([]database.DependencySource, error){ return []database.DependencySource{}, nil }
func (nd *NoopDatabase) BatchListDependencyTargetsOfDependency(ctx context.Context, dependencyIds []uuid.UUID) ([]database.DependencyTarget, error){ return []database.DependencyTarget{}, nil }

func (nd *NoopDatabase) WithTransaction(ctx context.Context, fn func(repo dbRepo.Database) error) error { return nil }
