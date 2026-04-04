-- +goose Up
-- +goose StatementBegin
CREATE TABLE statistik_ujian (
    id_statistik_ujian BIGSERIAL PRIMARY KEY,
    id_jadwal_ujian BIGINT NOT NULL REFERENCES jadwal_ujian (id_jadwal_ujian) ON DELETE CASCADE ON UPDATE CASCADE,
    nilai_tertinggi DECIMAL(5,2),
    nilai_terendah DECIMAL(5,2),
    nilai_rata_rata DECIMAL(5,2),
    total_peserta_ujian INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_statistik_ujian_id_jadwal_ujian
        UNIQUE (id_jadwal_ujian)
);
CREATE INDEX idx_statistik_ujian_id_jadwal_ujian
    ON statistik_ujian (id_jadwal_ujian);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE statistik_ujian;
-- +goose StatementEnd
