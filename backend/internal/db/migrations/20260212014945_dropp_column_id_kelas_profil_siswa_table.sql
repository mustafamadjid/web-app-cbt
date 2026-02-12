-- +goose Up
-- +goose StatementBegin
ALTER TABLE profil_siswa
    DROP COLUMN IF EXISTS id_kelas;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profil_siswa
    ADD COLUMN id_kelas BIGINT NOT NULL;    
-- +goose StatementEnd
