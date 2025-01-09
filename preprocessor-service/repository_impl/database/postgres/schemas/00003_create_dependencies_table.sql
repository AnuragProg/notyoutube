-- +goose Up
CREATE TABLE dependencies (
    id UUID NOT NULL,
    dag_id UUID NOT NULL REFERENCES dags(id),
    PRIMARY KEY (id, dag_id)
); 
CREATE TABLE dependency_sources(
    id UUID PRIMARY KEY,
    dag_id UUID NOT NULL,
    dependency_id UUID NOT NULL,
    source_id UUID NOT NULL,
    FOREIGN KEY (dag_id, dependency_id) REFERENCES dependencies(dag_id, id),
    FOREIGN KEY (source_id, dag_id) REFERENCES workers(id, dag_id)
);
CREATE TABLE dependency_targets(
    id UUID PRIMARY KEY,
    dag_id UUID NOT NULL,
    dependency_id UUID NOT NULL,
    target_id UUID NOT NULL,
    FOREIGN KEY (dag_id, dependency_id) REFERENCES dependencies(dag_id, id),
    FOREIGN KEY (target_id, dag_id) REFERENCES workers(id, dag_id)
);

-- +goose Down
DROP TABLE dependency_sources;
DROP TABLE dependency_targets;
DROP TABLE dependencies;
