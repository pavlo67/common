CREATE TABLE tagged (
  key      TEXT NOT NULL,
  id       TEXT NOT NULL,
  tag      TEXT NOT NULL,
  relation TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_tagged_uniq ON tagged(key, id, tag);

CREATE        INDEX idx_tagged      ON tagged(tag);

-------------------------

CREATE TABLE tags (
  tag         TEXT    NOT NULL,
  is_internal INTEGER NOT NULL,
  parted_size INTEGER NOT NULL
);

CREATE UNIQUE INDEX idx_tags_uniq ON tags(tag);

CREATE        INDEX idx_tags_int  ON tags(is_internal);
