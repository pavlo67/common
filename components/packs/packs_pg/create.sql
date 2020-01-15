CREATE TABLE packs (
    id                BIGSERIAL                PRIMARY KEY,
    identity_key      TEXT                     NOT NULL,
    address_from      TEXT                     NOT NULL,
    address_to        TEXT                     NOT NULL,
    options           TEXT                     NOT NULL,
    type_key          TEXT                     NOT NULL,
    content           TEXT                     NOT NULL,
    history           TEXT                     NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_packs_identity_key ON packs(identity_key);

CREATE INDEX idx_packs_type_key     ON packs(type_key);

CREATE INDEX idx_packs_created_at   ON packs(created_at);



