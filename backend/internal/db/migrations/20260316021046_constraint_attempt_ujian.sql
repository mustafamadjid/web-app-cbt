-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX uq_attempt_active
ON attempt_ujian (id_peserta_ujian)
WHERE status_attempt = 'in_progress';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uq_attempt_active;
-- +goose StatementEnd
