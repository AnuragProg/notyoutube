-- +goose Up
CREATE TYPE dag_status AS ENUM(
    'pending',
    'running',
    'completed',
    'failed'
);
CREATE TABLE dag_states(
    id UUID PRIMARY KEY,
    dag_id UUID NOT NULL,           -- dag id as stored in preprocessor service
    dag_status dag_status NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    failure_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE dag_states;
DROP TYPE dag_status;

