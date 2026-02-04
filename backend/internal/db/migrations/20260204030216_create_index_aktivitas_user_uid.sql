-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_aktivitas_user_id_pengguna
ON aktivitas_user(id_pengguna);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_aktivitas_user_id_pengguna;
-- +goose StatementEnd
