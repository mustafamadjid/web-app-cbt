-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS sessions_one_unrevoked_per_user
ON sessions(id_pengguna)
WHERE revoked_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS sessions_one_unrevoked_per_user;
-- +goose StatementEnd
