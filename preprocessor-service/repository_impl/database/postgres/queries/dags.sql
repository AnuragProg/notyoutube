-- name: CreateDAG :one
INSERT INTO dags (id)
VALUES ($1)
RETURNING *;

-- name: GetDAG :one
SELECT * FROM dags WHERE id = $1;
