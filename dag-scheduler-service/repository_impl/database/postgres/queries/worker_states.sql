-- name: CreateWorkerStates :batchexec
INSERT INTO worker_states(id, dag_id, worker_id, worker_status)
VALUES ($1, $2, $3, $4);

-- name: UpdateWorkerStateStartTime :exec
UPDATE worker_states
SET start_time=$2
WHERE id=$1;

-- name: UpdateWorkerStateEndTime :exec
UPDATE worker_states
SET end_time=$2
WHERE id=$1;

-- name: UpdateWorkerStateFailureReason :exec
UPDATE worker_states
SET failure_reason=$2
WHERE id=$1;

-- name: ListWorkerStatesOfWorker :many
SELECT *
FROM worker_states
WHERE worker_id=$1;

-- name: GetCurrentWorkerStateOfWorker :one
SELECT *
FROM worker_states
WHERE worker_id=$1
ORDER BY created_at DESC
LIMIT 1;
