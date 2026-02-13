-- +goose Up
-- +goose StatementBegin
ALTER TABLE ruang_ujian
    ADD COLUMN kode_ruang VARCHAR(100) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ruang_ujian
    DROP COLUMN IF EXISTS kode_ruang;
-- +goose StatementEnd
