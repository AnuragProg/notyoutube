-- name: CreateDAG :one
INSERT INTO dags (id, trace_id, created_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDAG :one
SELECT * FROM dags WHERE id = $1;
