-- +goose Up
-- +goose StatementBegin
ALTER TABLE isi_soal
    ADD COLUMN IF NOT EXISTS pertanyaan_content JSONB NULL;

ALTER TABLE opsi_pilihan_ganda
    ADD COLUMN IF NOT EXISTS isi_pilihan_content JSONB NULL;

ALTER TABLE import_soal_job
    ADD COLUMN IF NOT EXISTS warning_msg TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE import_soal_job
    DROP COLUMN IF EXISTS warning_msg;

ALTER TABLE opsi_pilihan_ganda
    DROP COLUMN IF EXISTS isi_pilihan_content;

ALTER TABLE isi_soal
    DROP COLUMN IF EXISTS pertanyaan_content;
-- +goose StatementEnd
