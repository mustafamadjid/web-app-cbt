-- +goose Up
-- +goose StatementBegin
ALTER TABLE profil_siswa
DROP CONSTRAINT IF EXISTS  uq_siswa_nisn;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profil_siswa
ADD CONSTRAINT uq_siswa_nisn UNIQUE (nisn);
-- +goose StatementEnd
