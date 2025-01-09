-- +goose Up
CREATE TYPE worker_type AS ENUM (
    'video-encoder',
    'ascii-encoder',
    'thumbnail-generator',
    'assembler',
    'video-extractor',
    'audio-extractor',
    'metadata-extractor'
);
CREATE TABLE workers (
    id UUID NOT NULL,
    dag_id UUID NOT NULL REFERENCES dags(id),
    name TEXT NOT NULL,
    description TEXT,
    worker_type worker_type NOT NULL,
    worker_config JSONB,
    PRIMARY KEY (id, dag_id)
);

-- +goose Down
DROP TABLE workers;
