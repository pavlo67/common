DROP TABLE IF EXISTS persons;


CREATE TABLE persons (
  id           INTEGER    PRIMARY KEY AUTOINCREMENT,
  issued_id    TEXT       NOT NULL,
  nickname     TEXT       NOT NULL,
  email        TEXT       NOT NULL,
  roles        TEXT       NOT NULL,
  creds        TEXT       NOT NULL,
  data         TEXT       NOT NULL,
  history      TEXT       NOT NULL,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE INDEX idx_persons_issued_id    ON persons(issued_id);
CREATE INDEX idx_persons_nickname     ON persons(nickname);
CREATE INDEX idx_persons_email        ON persons(email);


