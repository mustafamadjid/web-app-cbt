-- +goose Up
-- +goose StatementBegin
ALTER TABLE hasil_ujian
ADD COLUMN IF NOT EXISTS id_jadwal_ujian BIGINT NOT NULL REFERENCES jadwal_ujian(id_jadwal_ujian) ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE hasil_ujian DROP COLUMN IF EXISTS id_jadwal_ujian;
-- +goose StatementEnd
