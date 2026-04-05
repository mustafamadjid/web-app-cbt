-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS sessions_one_unrevoked_per_user;

CREATE UNIQUE INDEX IF NOT EXISTS sessions_one_unrevoked_non_admin_per_user
ON sessions(id_pengguna)
WHERE revoked_at IS NULL
  AND role = 'SISWA' OR role = 'GURU';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS sessions_one_unrevoked_non_admin_per_user;

CREATE UNIQUE INDEX IF NOT EXISTS sessions_one_unrevoked_per_user
ON sessions(id_pengguna)
WHERE revoked_at IS NULL;
-- +goose StatementEnd
