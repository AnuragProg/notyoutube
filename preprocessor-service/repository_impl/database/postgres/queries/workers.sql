
-- name: CreateWorkers :batchexec
INSERT INTO workers(id, dag_id, name, description, worker_type, worker_config)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListWorkersOfDAG :many
SELECT * FROM workers WHERE dag_id = $1;

-- name: GetWorkerById :one
SELECT *
FROM workers
WHERE id=$1
LIMIT 1;
