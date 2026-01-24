-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  id_pengguna BIGINT NOT NULL,
  expires_at timestamptz NOT NULL,
  revoked_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT fk_sessions_id_pengguna
    FOREIGN KEY (id_pengguna)
    REFERENCES pengguna(id_pengguna)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);
CREATE INDEX idx_sessions_id_pengguna ON sessions(id_pengguna);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS pgcrypto;

-- +goose StatementEnd
