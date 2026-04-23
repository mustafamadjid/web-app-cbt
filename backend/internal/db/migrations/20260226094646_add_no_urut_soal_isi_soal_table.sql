-- +goose Up
-- +goose StatementBegin
ALTER TABLE isi_soal ADD COLUMN IF NOT EXISTS no_urut_soal INTEGER DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE isi_soal DROP COLUMN IF EXISTS no_urut_soal;
-- +goose StatementEnd
