-- +goose Up
-- +goose StatementBegin
CREATE TABLE attempt_ujian (
    id_attempt BIGSERIAL PRIMARY KEY,
    id_peserta_ujian BIGINT NOT NULL,
    attempt_no BIGINT NOT NULL,
    status_attempt VARCHAR(20) NOT NULL DEFAULT 'in_progress', --in_progress, submitted, expired, cancelled
    waktu_mulai TIMESTAMPTZ DEFAULT NULL,
    waktu_submit TIMESTAMPTZ DEFAULT NULL,
    deadline_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_attempt_ujian_peserta_ujian
        FOREIGN KEY (id_peserta_ujian)
        REFERENCES peserta_ujian(id_peserta_ujian)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    
    CONSTRAINT chk_attempt_ujian_status
        CHECK (status_attempt IN ('in_progress', 'submitted', 'expired', 'cancelled'))
);

    CREATE UNIQUE INDEX uq_attempt_ujian_one_submitted_per_peserta
    ON attempt_ujian (id_peserta_ujian)
    WHERE status_attempt = 'submitted';

CREATE TABLE hasil_ujian (
    id_hasil BIGSERIAL PRIMARY KEY,
    id_attempt BIGINT NOT NULL,
    graded_by BIGINT DEFAULT NULL,
    nilai_akhir DECIMAL(5,2) DEFAULT NULL,
    passed BOOLEAN DEFAULT NULL,
    graded_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT fk_hasil_ujian_attempt
        FOREIGN KEY (id_attempt)
        REFERENCES attempt_ujian(id_attempt)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_hasil_ujian_pengguna
        FOREIGN KEY (graded_by)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    
    CONSTRAINT chk_hasil_ujian_nilai
        CHECK (nilai_akhir IS NULL OR (nilai_akhir >= 0 AND nilai_akhir <= 100))
);

CREATE INDEX idx_hasil_ujian_id_attempt ON hasil_ujian (id_attempt);
CREATE INDEX idx_hasil_ujian_graded_by ON hasil_ujian (graded_by);
CREATE INDEX idx_attempt_ujian_id_peserta_ujian ON attempt_ujian (id_peserta_ujian);

ALTER TABLE jawaban_ujian_siswa
    DROP COLUMN id_peserta_ujian;

ALTER TABLE jawaban_ujian_siswa
    ADD COLUMN id_attempt BIGINT NOT NULL,
    ADD CONSTRAINT fk_jawaban_ujian_siswa_attempt
        FOREIGN KEY (id_attempt)
        REFERENCES attempt_ujian(id_attempt)
        ON UPDATE CASCADE
        ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jawaban_ujian_siswa
    DROP CONSTRAINT fk_jawaban_ujian_siswa_attempt;

ALTER TABLE jawaban_ujian_siswa
    ADD COLUMN id_peserta_ujian BIGINT NOT NULL,
    ADD CONSTRAINT fk_jawaban_ujian_siswa_peserta_ujian
        FOREIGN KEY (id_peserta_ujian)
        REFERENCES peserta_ujian(id_peserta_ujian)
        ON UPDATE CASCADE
        ON DELETE CASCADE;

DROP TABLE attempt_ujian;
DROP TABLE hasil_ujian;
-- +goose StatementEnd
