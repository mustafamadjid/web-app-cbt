-- +goose Up
-- +goose StatementBegin
ALTER TABLE bank_soal
    ADD COLUMN IF NOT EXISTS id_bank_soal_version_aktif BIGINT NULL;

ALTER TABLE bank_soal
    ADD CONSTRAINT fk_bank_soal_version_aktif
    FOREIGN KEY (id_bank_soal_version_aktif)
    REFERENCES bank_soal_version(id_bank_soal_version)
    ON UPDATE CASCADE
    ON DELETE SET NULL;

ALTER TABLE isi_soal
ADD COLUMN IF NOT EXISTS id_bank_soal_version BIGINT NULL;

ALTER TABLE isi_soal
    ADD CONSTRAINT fk_isi_soal_bank_soal_version
    FOREIGN KEY (id_bank_soal_version)
    REFERENCES bank_soal_version(id_bank_soal_version)
    ON UPDATE CASCADE
    ON DELETE CASCADE;

ALTER TABLE isi_soal
    DROP CONSTRAINT fk_soal_bank_soal;

ALTER TABLE isi_soal
    DROP COLUMN id_bank_soal;
CREATE INDEX IF NOT EXISTS idx_isi_soal_id_bank_soal_version
    ON isi_soal(id_bank_soal_version);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_isi_soal_id_bank_soal_version;

ALTER TABLE isi_soal
DROP CONSTRAINT IF EXISTS fk_isi_soal_bank_soal_version;

ALTER TABLE isi_soal
DROP COLUMN IF EXISTS id_bank_soal_version;

ALTER TABLE bank_soal
DROP CONSTRAINT IF EXISTS fk_bank_soal_version_aktif;

ALTER TABLE bank_soal
DROP COLUMN IF EXISTS id_bank_soal_version_aktif;
-- +goose StatementEnd
