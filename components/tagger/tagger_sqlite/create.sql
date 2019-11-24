CREATE TABLE tags (
  key TEXT NOT NULL,
  id  TEXT NOT NULL,
  tag TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_data_source ON tags(key, id, tag);


