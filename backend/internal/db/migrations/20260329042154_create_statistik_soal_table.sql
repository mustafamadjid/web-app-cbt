-- +goose Up
-- +goose StatementBegin
CREATE TABLE statistik_soal (
    id_statistik_soal BIGSERIAL PRIMARY KEY,
    id_soal BIGINT NOT NULL REFERENCES isi_soal(id_soal) ON DELETE CASCADE ON UPDATE CASCADE,
    id_ujian BIGINT NOT NULL REFERENCES ujian(id_ujian) ON DELETE CASCADE ON UPDATE CASCADE,
    jumlah_soal INTEGER NOT NULL,
    jumlah_jawaban_benar INTEGER NOT NULL,
    jumlah_jawaban_salah INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_statistik_soal_id_ujian_id_soal
        UNIQUE (id_ujian, id_soal)
);

CREATE INDEX idx_statistik_soal_id_ujian
    ON statistik_soal (id_ujian);

CREATE INDEX idx_statistik_soal_id_soal
    ON statistik_soal (id_soal);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS statistik_soal;
-- +goose StatementEnd
