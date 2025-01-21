-- +goose Up
CREATE TYPE worker_status AS ENUM(
    'pending',
    'running',
    'completed',
    'failed'
);
CREATE TABLE worker_states(
    id UUID PRIMARY KEY,
    dag_id UUID NOT NULL,               -- dag id as stored in preprocessor service
    worker_id UUID NOT NULL,            -- worker id as stored in preprocessor service
    worker_status worker_status NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    failure_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
); 

-- +goose Down
DROP TABLE worker_states;
DROP TYPE worker_status;
