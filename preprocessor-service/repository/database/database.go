package database

import (
	"context"
	"io"

	"github.com/anuragprog/notyoutube/preprocessor-service/types/database"
	"github.com/google/uuid"
)

type Database interface {

	// to make sure the handler closes the connection properly
	io.Closer

	CreateDAG(ctx context.Context, dag database.Dag) error
	GetDAG(ctx context.Context, id uuid.UUID) (database.Dag, error)
	
	CreateWorkers(ctx context.Context, workers []database.Worker) error
	ListWorkersOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.Worker, error)
	GetWorkerById(ctx context.Context, workerId uuid.UUID) (database.Worker, error)

	CreateDependencies(ctx context.Context, dependencies []database.Dependency) error
	CreateDependencySources(ctx context.Context, sources []database.DependencySource) error
	CreateDependencyTargets(ctx context.Context, targets []database.DependencyTarget) error
	ListDependenciesOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.Dependency, error)
	ListDependencySourcesOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.DependencySource, error)
	ListDependencyTargetsOfDAG(ctx context.Context, dagId uuid.UUID) ([]database.DependencyTarget, error)
	ListDependencySourcesWhereWorkerIsSource(ctx context.Context, workerId uuid.UUID) ([]database.DependencySource, error) 
	BatchListDependencySourcesOfDependency(ctx context.Context, dependencyIds []uuid.UUID) ([]database.DependencySource, error)
	BatchListDependencyTargetsOfDependency(ctx context.Context, dependencyIds []uuid.UUID) ([]database.DependencyTarget, error)

	WithTransaction(ctx context.Context, fn func(repo Database) error) error
}
