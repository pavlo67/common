CREATE TABLE tags (
  id           INTEGER             PRIMARY KEY AUTOINCREMENT,
  saved_at     TIMESTAMP           DEFAULT CURRENT_TIMESTAMP,
  source_id    TEXT       NOT NULL,
  source_key   TEXT       NOT NULL,
  origin       TEXT       NOT NULL,
  source_time  TIMESTAMP,
  source_url   TEXT,
  types        TEXT       NOT NULL,
  title        TEXT       NOT NULL,
  summary      TEXT,
  details      TEXT,
  href         TEXT,
  embedded     TEXT,
  tags         TEXT,
  indexes      TEXT
);

CREATE INDEX idx_data_source     ON data(source_id, source_key);

CREATE INDEX idx_data_source_url ON data(source_url);

CREATE INDEX idx_data_types      ON data(types, saved_at);


