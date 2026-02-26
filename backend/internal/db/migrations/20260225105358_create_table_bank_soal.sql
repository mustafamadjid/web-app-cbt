-- +goose Up
-- +goose StatementBegin


CREATE TABLE IF NOT EXISTS bank_soal (
    id_bank_soal BIGSERIAL PRIMARY KEY,
    id_mapel BIGINT NOT NULL,
    id_kelas BIGINT NOT NULL,
    id_pengguna BIGINT NOT NULL,
    nama_bank_soal VARCHAR(100) NOT NULL,
    deskripsi TEXT NOT NULL,
    materi VARCHAR(100) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_bank_soal_mapel
        FOREIGN KEY (id_mapel)
        REFERENCES mata_pelajaran(id_mapel)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_bank_soal_kelas
        FOREIGN KEY (id_kelas)
        REFERENCES kelas(id_kelas)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_bank_soal_pengguna
        FOREIGN KEY (id_pengguna)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS isi_soal(
    id_soal BIGSERIAL PRIMARY KEY,
    id_bank_soal BIGINT NOT NULL,
    tipe_soal VARCHAR(20) NOT NULL,
    pertanyaan TEXT NOT NULL,
    gambar VARCHAR(255) DEFAULT NULL,
    bobot_soal INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_soal_bank_soal
        FOREIGN KEY (id_bank_soal)
        REFERENCES bank_soal(id_bank_soal)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS opsi_pilihan_ganda(
    id_pilihan_ganda BIGSERIAL PRIMARY KEY,
    id_soal BIGINT NOT NULL,
    isi_pilihan TEXT NOT NULL,
    is_benar BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_pilihan_ganda_soal
        FOREIGN KEY (id_soal)
        REFERENCES isi_soal(id_soal)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_bank_soal_id_mapel ON bank_soal(id_mapel);
CREATE INDEX IF NOT EXISTS idx_bank_soal_id_kelas ON bank_soal(id_kelas);
CREATE INDEX IF NOT EXISTS idx_bank_soal_id_pengguna ON bank_soal(id_pengguna);
CREATE INDEX IF NOT EXISTS idx_isi_soal_id_bank_soal ON isi_soal(id_bank_soal);
CREATE INDEX IF NOT EXISTS idx_pilihan_ganda_id_soal ON opsi_pilihan_ganda(id_soal);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS opsi_pilihan_ganda;
DROP TABLE IF EXISTS isi_soal;
DROP TABLE IF EXISTS bank_soal;
-- +goose StatementEnd
