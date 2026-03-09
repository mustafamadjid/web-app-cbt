-- +goose Up
-- +goose StatementBegin
ALTER TABLE peserta_ujian
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
