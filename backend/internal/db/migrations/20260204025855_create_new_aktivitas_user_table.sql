-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS aktivitas_user(
    id_aktivitas UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_pengguna  BIGINT NOT NULL,
    action      VARCHAR(20) NOT NULL,
    description  VARCHAR(255) NOT NULL,
    ip_address   VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),


    CONSTRAINT fk_aktivitas_user_id_pengguna
        FOREIGN KEY (id_pengguna)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS aktivitas_user;
-- +goose StatementEnd
