-- +goose Up
-- +goose StatementBegin

-- profil_siswa: pengguna -> cascade
ALTER TABLE profil_siswa
  DROP CONSTRAINT fk_siswa_pengguna;

ALTER TABLE profil_siswa
  ADD CONSTRAINT fk_siswa_pengguna
  FOREIGN KEY (id_pengguna)
  REFERENCES pengguna(id_pengguna)
  ON UPDATE CASCADE
  ON DELETE CASCADE;

-- profil_siswa: kelas -> cascade
ALTER TABLE profil_siswa
  DROP CONSTRAINT fk_siswa_kelas;

ALTER TABLE profil_siswa
  ADD CONSTRAINT fk_siswa_kelas
  FOREIGN KEY (id_kelas)
  REFERENCES kelas(id_kelas)
  ON UPDATE CASCADE
  ON DELETE CASCADE;

-- profil_guru: pengguna -> cascade (hapus kalau tidak diinginkan)
ALTER TABLE profil_guru
  DROP CONSTRAINT fk_guru_pengguna;

ALTER TABLE profil_guru
  ADD CONSTRAINT fk_guru_pengguna
  FOREIGN KEY (id_pengguna)
  REFERENCES pengguna(id_pengguna)
  ON UPDATE CASCADE
  ON DELETE CASCADE;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- profil_siswa: pengguna -> restrict
ALTER TABLE profil_siswa
  DROP CONSTRAINT fk_siswa_pengguna;

ALTER TABLE profil_siswa
  ADD CONSTRAINT fk_siswa_pengguna
  FOREIGN KEY (id_pengguna)
  REFERENCES pengguna(id_pengguna)
  ON UPDATE CASCADE
  ON DELETE RESTRICT;

-- profil_siswa: kelas -> restrict
ALTER TABLE profil_siswa
  DROP CONSTRAINT fk_siswa_kelas;

ALTER TABLE profil_siswa
  ADD CONSTRAINT fk_siswa_kelas
  FOREIGN KEY (id_kelas)
  REFERENCES kelas(id_kelas)
  ON UPDATE CASCADE
  ON DELETE RESTRICT;

-- profil_guru: pengguna -> restrict
ALTER TABLE profil_guru
  DROP CONSTRAINT fk_guru_pengguna;

ALTER TABLE profil_guru
  ADD CONSTRAINT fk_guru_pengguna
  FOREIGN KEY (id_pengguna)
  REFERENCES pengguna(id_pengguna)
  ON UPDATE CASCADE
  ON DELETE RESTRICT;

-- +goose StatementEnd
