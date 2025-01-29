
-- name: CreateDependencies :batchexec
INSERT INTO dependencies(id, dag_id)
VALUES ($1, $2);

-- name: CreateDependencySources :batchexec
INSERT INTO dependency_sources(id, dag_id, dependency_id, source_id)
VALUES ($1, $2, $3, $4);

-- name: CreateDependencyTargets :batchexec
INSERT INTO dependency_targets(id, dag_id, dependency_id, target_id)
VALUES ($1, $2, $3, $4);

-- name: ListDependenciesOfDAG :many
SELECT * FROM dependencies WHERE dag_id = $1;

-- name: ListDependencySourcesOfDependency :many
SELECT * FROM dependency_sources WHERE dependency_id = $1;

-- name: ListDependencyTargetsOfDependency :many
SELECT * FROM dependency_targets WHERE dependency_id = $1;

-- name: ListDependencySourcesWhereWorkerIsSource :many
SELECT * 
FROM dependency_sources
WHERE source_id=sqlc.arg(worker_id);

-- name: BatchListDependencySourcesOfDependency :batchmany
SELECT * 
FROM dependency_sources
WHERE dependency_id=$1;

-- name: BatchListDependencyTargetsOfDependency :batchmany
SELECT *
FROM dependency_targets
WHERE dependency_id=$1;
