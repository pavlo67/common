CREATE TABLE data (
  id           INTEGER             PRIMARY KEY AUTOINCREMENT,
  export_id    TEXT       NOT NULL,
  url          TEXT       NOT NULL,
  type         TEXT       NOT NULL,
  title        TEXT       NOT NULL,
  summary      TEXT       NOT NULL,
  embedded     TEXT       NOT NULL,
  tags         TEXT       NOT NULL,
  details      TEXT       NOT NULL,
  source       TEXT       NOT NULL,
  source_key   TEXT       NOT NULL,
  source_time  TIMESTAMP,
  source_data  TEXT       NOT NULL,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX idx_data_source     ON data(source, source_key);

CREATE INDEX idx_data_title      ON data(`type`, title);

--

CREATE TABLE flow (
  id           INTEGER             PRIMARY KEY AUTOINCREMENT,
  export_id    TEXT       NOT NULL,
  url          TEXT       NOT NULL,
  type         TEXT       NOT NULL,
  title        TEXT       NOT NULL,
  summary      TEXT       NOT NULL,
  embedded     TEXT       NOT NULL,
  tags         TEXT       NOT NULL,
  details      TEXT       NOT NULL,
  source       TEXT       NOT NULL,
  source_key   TEXT       NOT NULL,
  source_time  TIMESTAMP,
  source_data  TEXT       NOT NULL,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX idx_flow_source     ON flow(source, source_key);

CREATE INDEX idx_flow_title      ON flow(`type`, title);
