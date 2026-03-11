-- +goose Up
-- +goose StatementBegin
ALTER TABLE jadwal_ujian
ADD COLUMN status_ujian VARCHAR(100) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jadwal_ujian
DROP COLUMN IF EXISTS status_ujian;
-- +goose StatementEnd
