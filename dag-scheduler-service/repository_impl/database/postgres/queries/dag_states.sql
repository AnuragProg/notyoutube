-- name: CreateDAGState :one
INSERT INTO dag_states (
    id, dag_id, dag_status
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateDAGStateStartTime :exec
UPDATE dag_states
SET start_time=$2
WHERE id=$1;

-- name: UpdateDAGStateEndTime :exec
UPDATE dag_states
SET end_time=$2
WHERE id=$1;

-- name: UpdateDAGStateFailureReason :exec
UPDATE dag_states
SET failure_reason=$2
WHERE id=$1;

-- name: GetCurrentDAGState :one
SELECT *
FROM dag_states
WHERE dag_id=$1
ORDER BY created_at DESC
LIMIT 1;

-- name: ListDAGStates :many
SELECT * 
FROM dag_states 
WHERE dag_id = $1;

