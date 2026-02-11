-- +goose Up
-- +goose StatementBegin
ALTER TABLE nama_kelas
  DROP CONSTRAINT fk_nama_kelas_kelas;

ALTER TABLE nama_kelas
  ADD CONSTRAINT fk_nama_kelas_kelas
  FOREIGN KEY (id_kelas)
  REFERENCES kelas(id_kelas)
  ON UPDATE CASCADE
  ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE nama_kelas
  DROP CONSTRAINT fk_nama_kelas_kelas;

ALTER TABLE nama_kelas
  ADD CONSTRAINT fk_nama_kelas_kelas
  FOREIGN KEY (id_kelas)
  REFERENCES kelas(id_kelas)
  ON UPDATE CASCADE
  ON DELETE RESTRICT;

-- +goose StatementEnd
