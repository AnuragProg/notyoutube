-- +goose Up
CREATE TABLE dags (
    id UUID PRIMARY KEY,
    trace_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE dags;
