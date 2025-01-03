// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: dependencies.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listDependenciesOfDAG = `-- name: ListDependenciesOfDAG :many
SELECT id, dag_id FROM dependencies WHERE dag_id = $1
`

func (q *Queries) ListDependenciesOfDAG(ctx context.Context, dagID pgtype.UUID) ([]Dependency, error) {
	rows, err := q.db.Query(ctx, listDependenciesOfDAG, dagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Dependency
	for rows.Next() {
		var i Dependency
		if err := rows.Scan(&i.ID, &i.DagID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDependencySourcesOfDependency = `-- name: ListDependencySourcesOfDependency :many
SELECT id, dag_id, dependency_id, source_id FROM dependency_sources WHERE dependency_id = $1
`

func (q *Queries) ListDependencySourcesOfDependency(ctx context.Context, dependencyID pgtype.UUID) ([]DependencySource, error) {
	rows, err := q.db.Query(ctx, listDependencySourcesOfDependency, dependencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DependencySource
	for rows.Next() {
		var i DependencySource
		if err := rows.Scan(
			&i.ID,
			&i.DagID,
			&i.DependencyID,
			&i.SourceID,
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

const listDependencyTargetsOfDependency = `-- name: ListDependencyTargetsOfDependency :many
SELECT id, dag_id, dependency_id, target_id FROM dependency_targets WHERE dependency_id = $1
`

func (q *Queries) ListDependencyTargetsOfDependency(ctx context.Context, dependencyID pgtype.UUID) ([]DependencyTarget, error) {
	rows, err := q.db.Query(ctx, listDependencyTargetsOfDependency, dependencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DependencyTarget
	for rows.Next() {
		var i DependencyTarget
		if err := rows.Scan(
			&i.ID,
			&i.DagID,
			&i.DependencyID,
			&i.TargetID,
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
