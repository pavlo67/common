// sources ---------------------------------------------------------

CREATE TABLE sources (
  id           BIGINT UNSIGNED NOT NULL PRIMARY KEY,
  saved_at     TIMESTAMP                DEFAULT CURRENT_TIMESTAMP,
  title        TEXT            NOT NULL,
  url          TEXT            NOT NULL,
  tags         TEXT
);

CREATE INDEX idx_sources_title    ON sources(title);

CREATE INDEX idx_sources_url      ON sources(url);

CREATE INDEX idx_sources_saved_at ON sources(saved_at);


// flow --------------------------------------------------

CREATE TABLE flow (
  id           BIGINT UNSIGNED NOT NULL PRIMARY KEY,
  saved_at     TIMESTAMP                DEFAULT CURRENT_TIMESTAMP,
  source_id    BIGINT UNSIGNED NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
  source_key   TEXT            NOT NULL,
  origin       TEXT            NOT NULL,
  source_time  TIMESTAMP,
  source_url   TEXT,
  title        TEXT            NOT NULL,
  summary      TEXT,
  details      TEXT,
  href         TEXT,
  embedded     TEXT,
  tags         TEXT
);

CREATE INDEX idx_flow_source   ON flow(source_id, source_key);

CREATE INDEX idx_flow_saved_at ON flow(saved_at);


// tags -------------------------------------------------

CREATE TABLE tags (
  tag          TEXT            NOT NULL,
  saved_at     TIMESTAMP                DEFAULT CURRENT_TIMESTAMP,
  source_id    BIGINT UNSIGNED NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
  flow_id      BIGINT UNSIGNED NOT NULL REFERENCES flow(id)    ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_tags_tag       ON tags(tag, source_id);

CREATE        INDEX idx_tags_saved_at  ON tags(saved_at);

CREATE        INDEX idx_tags_flow_id   ON tags(flow_id);

CREATE        INDEX idx_tags_source_id ON tags(source_id);
