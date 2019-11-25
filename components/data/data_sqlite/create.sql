CREATE TABLE data (
  id           INTEGER             PRIMARY KEY AUTOINCREMENT,
  title        TEXT       NOT NULL,
  summary      TEXT,
  url          TEXT,
  embedded     TEXT,
  tags         TEXT,
  details      TEXT,
  source       TEXT,
  source_key   TEXT,
  source_time  TIMESTAMP,
  source_data  TEXT,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX idx_data_source     ON data(source, source_key);

