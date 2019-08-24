// flow --------------------------------------------------

CREATE TABLE flow (
  id           INTEGER             PRIMARY KEY AUTOINCREMENT,
  saved_at     TIMESTAMP           DEFAULT CURRENT_TIMESTAMP,
  source_id    TEXT       NOT NULL,
  source_key   TEXT       NOT NULL,
  origin       TEXT       NOT NULL,
  source_time  TIMESTAMP,
  source_url   TEXT,
  title        TEXT       NOT NULL,
  summary      TEXT,
  details      TEXT,
  href         TEXT,
  embedded     TEXT,
  tags         TEXT
);

CREATE INDEX idx_flow_source     ON flow(source_id, source_key);

CREATE INDEX idx_flow_source_url ON flow(source_url);

CREATE INDEX idx_flow_saved_at   ON flow(saved_at);


// tags -------------------------------------------------

CREATE TABLE tags (
  tag      TEXT      NOT NULL,
  saved_at TIMESTAMP          DEFAULT CURRENT_TIMESTAMP,
  flow_id  INTEGER   NOT NULL REFERENCES flow(id)    ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_tags_tag      ON tags(tag, flow_id);

CREATE        INDEX idx_tags_flow_id  ON tags(flow_id);

CREATE        INDEX idx_tags_saved_at ON tags(saved_at);

