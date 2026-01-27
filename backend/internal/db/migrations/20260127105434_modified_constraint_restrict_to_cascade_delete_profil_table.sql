-- +goose Up
-- +goose StatementBegin
-- profil_siswa: pengguna -> cascade
ALTER TABLE profil_siswa
DROP CONSTRAINT IF EXISTS fk_siswa_pengguna;

ALTER TABLE profil_siswa
ADD CONSTRAINT fk_siswa_pengguna
FOREIGN KEY (id_pengguna)
REFERENCES pengguna(id_pengguna)
ON UPDATE CASCADE
ON DELETE CASCADE;

-- profil_guru: pengguna -> cascade
ALTER TABLE profil_guru
DROP CONSTRAINT IF EXISTS fk_guru_pengguna;

ALTER TABLE profil_guru
ADD CONSTRAINT fk_guru_pengguna
FOREIGN KEY (id_pengguna)
REFERENCES pengguna(id_pengguna)
ON UPDATE CASCADE
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profil_siswa
DROP CONSTRAINT IF EXISTS fk_siswa_pengguna;

ALTER TABLE profil_siswa
ADD CONSTRAINT fk_siswa_pengguna
FOREIGN KEY (id_pengguna)
REFERENCES pengguna(id_pengguna)
ON UPDATE CASCADE
ON DELETE RESTRICT;

ALTER TABLE profil_guru
DROP CONSTRAINT IF EXISTS fk_guru_pengguna;

ALTER TABLE profil_guru
ADD CONSTRAINT fk_guru_pengguna
FOREIGN KEY (id_pengguna)
REFERENCES pengguna(id_pengguna)
ON UPDATE CASCADE
ON DELETE RESTRICT;
-- +goose StatementEnd
