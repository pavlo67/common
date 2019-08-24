// flow --------------------------------------------------

CREATE TABLE datas (
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

CREATE INDEX idx_datas_source     ON datas(source_id, source_key);

CREATE INDEX idx_datas_source_url ON datas(source_url);

CREATE INDEX idx_datas_types      ON datas(types, saved_at);


