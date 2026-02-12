-- +goose Up
-- +goose StatementBegin
ALTER TABLE profil_siswa
    DROP CONSTRAINT IF EXISTS fk_siswa_kelas;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profil_siswa
    ADD CONSTRAINT fk_siswa_kelas
    FOREIGN KEY (id_kelas)
    REFERENCES kelas(id_kelas)
    ON UPDATE CASCADE
    ON DELETE RESTRICT;
-- +goose StatementEnd
