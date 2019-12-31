CREATE TABLE storage (
  id           INTEGER    PRIMARY KEY AUTOINCREMENT,
  data_key     TEXT       ,
  url          TEXT       NOT NULL,
  type         TEXT       NOT NULL,
  title        TEXT       NOT NULL,
  summary      TEXT       NOT NULL,
  embedded     TEXT       NOT NULL,
  tags         TEXT       NOT NULL,
  details      TEXT       NOT NULL,
  history      TEXT       ,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX idx_storage_key   ON storage(data_key);

CREATE INDEX idx_storage_title ON storage(`type`, title);

--

CREATE TABLE flow (
  id           INTEGER    PRIMARY KEY AUTOINCREMENT,
  data_key     TEXT       ,
  url          TEXT       NOT NULL,
  type         TEXT       NOT NULL,
  title        TEXT       NOT NULL,
  summary      TEXT       NOT NULL,
  embedded     TEXT       NOT NULL,
  tags         TEXT       NOT NULL,
  details      TEXT       NOT NULL,
  history      TEXT       ,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX idx_flow_key   ON flow(data_key);

CREATE INDEX idx_flow_title ON flow(`type`, title);
