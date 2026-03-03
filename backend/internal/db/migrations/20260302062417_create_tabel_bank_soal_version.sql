-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bank_soal_version (
    id_bank_soal_version BIGSERIAL PRIMARY KEY,
    id_bank_soal BIGINT NOT NULL,
    version_no INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft', -- draft|published|archived
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by BIGINT NULL,

    CONSTRAINT fk_bank_soal_version_bank
        FOREIGN KEY (id_bank_soal)
        REFERENCES bank_soal(id_bank_soal)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT uq_bank_soal_version UNIQUE (id_bank_soal, version_no)
);

CREATE INDEX IF NOT EXISTS idx_bank_soal_version_bank
ON bank_soal_version(id_bank_soal);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bank_soal_version;
-- +goose StatementEnd
