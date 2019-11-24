// flow --------------------------------------------------

CREATE TABLE data (
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



{
    sql: {
        data: {
            fields: [
                ["id",       "INTEGER", "true", "", "PRIMARY KEY AUTOINCREMENT"],
                ["saved_at", "TIMESTAMP", "false", "CURRENT_TIMESTAMP", ""],

                ["source_id",   "TEXT", "false", "", ""],
                ["source_key",  "TEXT", "false", "", ""],
                ["origin",      "TEXT", "false", "", ""],
                ["source_time", "TIMESTAMP", "true", "", ""],
                ["source_url",  "TEXT", "true", "", ""],

                ["types",    "TEXT", "false", "", ""],
                ["title",    "TEXT", "false", "", ""],
                ["summary",  "TEXT", "true", "", ""],
                ["details",  "TEXT", "true", "", ""],
                ["href",     "TEXT", "true", "", ""],
                ["embedded", "TEXT", "true", "", ""],

                ["tags",    "TEXT", "true", "", ""],
                ["indexes", "TEXT", "true", "", ""],
            ],
            indexes: [
                {name: "source_url", type: "", fields: ["source_url"]},
                {name: "source",     type: "", fields: ["source_id", "source_key"]},
                {name: "types",      type: "", fields: ["types", "saved_at"]},
            ],
        },
    },
}

