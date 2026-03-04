-- +goose Up
-- +goose StatementBegin
ALTER TABLE jadwal_ujian
  ADD COLUMN IF NOT EXISTS token VARCHAR(100) NOT NULL,
  ADD CONSTRAINT ck_jadwal_ujian_token_uppercase
CHECK (token = upper(token));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jadwal_ujian
DROP CONSTRAINT IF EXISTS ck_jadwal_ujian_token_uppercase;

ALTER TABLE jadwal_ujian
DROP COLUMN IF EXISTS token;
-- +goose StatementEnd
