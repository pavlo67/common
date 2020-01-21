CREATE TABLE tasks (
    id                BIGSERIAL                PRIMARY KEY,
    worker_type       VARCHAR(255)             NOT NULL,
    params            TEXT                     ,
    status            TEXT                     NOT NULL,
    results           TEXT                     ,
    history           TEXT                     ,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_type       ON tasks(worker_type);

CREATE INDEX idx_tasks_created_at ON tasks(created_at);


