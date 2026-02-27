-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS import_soal_job (
    id_job       BIGSERIAL PRIMARY KEY,
    id_bank_soal BIGINT NOT NULL,
    id_pengguna  BIGINT NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending',
    file_path    TEXT NOT NULL,
    error_msg    TEXT,
    total_soal   INT DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_job_bank_soal
        FOREIGN KEY (id_bank_soal)
        REFERENCES bank_soal(id_bank_soal)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_job_pengguna
        FOREIGN KEY (id_pengguna)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_import_job_status ON import_soal_job(status);
CREATE INDEX IF NOT EXISTS idx_import_job_bank_soal ON import_soal_job(id_bank_soal);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS import_soal_job;

-- +goose StatementEnd
