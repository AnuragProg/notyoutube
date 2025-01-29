// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workers.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getWorkerById = `-- name: GetWorkerById :one
SELECT id, dag_id, name, description, worker_type, worker_config
FROM workers
WHERE id=$1
LIMIT 1
`

func (q *Queries) GetWorkerById(ctx context.Context, id pgtype.UUID) (Worker, error) {
	row := q.db.QueryRow(ctx, getWorkerById, id)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.DagID,
		&i.Name,
		&i.Description,
		&i.WorkerType,
		&i.WorkerConfig,
	)
	return i, err
}

const listWorkersOfDAG = `-- name: ListWorkersOfDAG :many
SELECT id, dag_id, name, description, worker_type, worker_config FROM workers WHERE dag_id = $1
`

func (q *Queries) ListWorkersOfDAG(ctx context.Context, dagID pgtype.UUID) ([]Worker, error) {
	rows, err := q.db.Query(ctx, listWorkersOfDAG, dagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Worker
	for rows.Next() {
		var i Worker
		if err := rows.Scan(
			&i.ID,
			&i.DagID,
			&i.Name,
			&i.Description,
			&i.WorkerType,
			&i.WorkerConfig,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
