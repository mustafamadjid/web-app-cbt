-- +goose Up
-- +goose StatementBegin
ALTER TABLE jadwal_ujian
    ADD COLUMN IF NOT EXISTS id_pengawas BIGINT,
    ADD CONSTRAINT fk_jadwal_ujian_pengawas
        FOREIGN KEY (id_pengawas)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jadwal_ujian
DROP CONSTRAINT IF EXISTS fk_jadwal_ujian_pengawas;

ALTER TABLE jadwal_ujian
DROP COLUMN IF EXISTS id_pengawas;
-- +goose StatementEnd
