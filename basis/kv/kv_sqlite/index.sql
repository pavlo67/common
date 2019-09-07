CREATE TABLE kv (
  key          TEXT      NOT NULL PRIMARY KEY,
  value        TEXT,
  saved_at     TIMESTAMP          DEFAULT CURRENT_TIMESTAMP
)

