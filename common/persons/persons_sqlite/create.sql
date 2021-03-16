DROP TABLE IF EXISTS persons;

CREATE TABLE persons (
  id           INTEGER    PRIMARY KEY AUTOINCREMENT,
  nickname     TEXT       NOT NULL,
  email        TEXT               ,
  roles        TEXT       NOT NULL,
  creds        TEXT       NOT NULL,

  label        TEXT       NOT NULL,
  info         TEXT       NOT NULL,
  tags         TEXT       NOT NULL,
  urn          TEXT               ,
  pack_urn     TEXT       NOT NULL,
  owner_nss    TEXT       NOT NULL,
  viewer_nss   TEXT       NOT NULL,
  history      TEXT       NOT NULL,
  created_at   TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP
);

CREATE UNIQUE INDEX idx_persons_email    ON persons(email)    WHERE email IS NOT NULL;
CREATE UNIQUE INDEX idx_persons_urn      ON persons(urn)      WHERE urn   IS NOT NULL;
CREATE        INDEX idx_persons_nickname ON persons(nickname);
CREATE        INDEX idx_persons_pack_urn ON persons(pack_urn);
