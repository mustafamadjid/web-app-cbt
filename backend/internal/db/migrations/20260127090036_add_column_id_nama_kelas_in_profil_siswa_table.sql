-- +goose Up
-- +goose StatementBegin
ALTER TABLE profil_siswa 
    ADD COLUMN id_nama_kelas BIGINT NOT NULL
    REFERENCES nama_kelas(id_nama_kelas)
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profil_siswa
    DROP COLUMN  IF EXISTS id_nama_kelas;

-- +goose StatementEnd
