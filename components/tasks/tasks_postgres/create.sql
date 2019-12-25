CREATE TABLE tasks (
    id                BIGSERIAL                PRIMARY KEY,
    worker_type       VARCHAR(255)             NOT NULL,
    params            TEXT                     NOT NULL,
    status            TEXT                     NOT NULL,
    results           TEXT                     NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tasks_type       ON tasks(worker_type);

CREATE INDEX idx_tasks_created_at ON tasks(created_at);


