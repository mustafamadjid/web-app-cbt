-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ujian (
    id_ujian BIGSERIAL PRIMARY KEY,
    id_bank_soal BIGINT NOT NULL,
    id_kelas BIGINT NOT NULL,
    id_nama_kelas BIGINT DEFAULT NULL,
    id_guru BIGINT NOT NULL,
    nama_ujian VARCHAR(100) NOT NULL,
    deskripsi_ujian TEXT DEFAULT NULL,
    acak_soal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT fk_ujian_bank_soal
        FOREIGN KEY (id_bank_soal)
        REFERENCES bank_soal(id_bank_soal)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_ujian_guru
        FOREIGN KEY (id_guru)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_ujian_kelas
        FOREIGN KEY (id_kelas)
        REFERENCES kelas(id_kelas)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_ujian_nama_kelas
        FOREIGN KEY (id_nama_kelas)
        REFERENCES nama_kelas(id_nama_kelas)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
    
);


CREATE TABLE IF NOT EXISTS jadwal_ujian (
    id_jadwal_ujian BIGSERIAL PRIMARY KEY,
    id_ujian BIGINT NOT NULL,
    id_sesi BIGINT NOT NULL,
    id_ruangan BIGINT NOT NULL,
    tanggal_ujian DATE NOT NULL,
    waktu_mulai TIMESTAMPTZ NOT NULL,
    waktu_selesai TIMESTAMPTZ NOT NULL,
    status_ujian VARCHAR(15) NOT NULL DEFAULT 'BELUM_MULAI', -- BELUM_MULAI | MULAI | SELESAI | DIBATALKAN
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT fk_jadwal_ujian_ujian
        FOREIGN KEY (id_ujian)
        REFERENCES ujian(id_ujian)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_jadwal_ujian_sesi
        FOREIGN KEY (id_sesi)
        REFERENCES sesi_ujian(id_sesi)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_jadwal_ujian_ruangan
        FOREIGN KEY (id_ruangan)
        REFERENCES ruang_ujian(id_ruangan)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

        CONSTRAINT ck_status_ujian
            CHECK (status_ujian IN ('BELUM_MULAI', 'MULAI', 'SELESAI', 'DIBATALKAN')),

    CONSTRAINT ck_waktu CHECK (waktu_mulai < waktu_selesai)
);


CREATE TABLE IF NOT EXISTS peserta_ujian (
    id_peserta_ujian BIGSERIAL PRIMARY KEY,
    id_jadwal_ujian BIGINT NOT NULL,
    id_siswa BIGINT NOT NULL,
    waktu_mulai TIMESTAMPTZ,
    waktu_submit TIMESTAMPTZ,
    nilai_ujian DECIMAL(5,2) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT fk_peserta_ujian_jadwal_ujian
        FOREIGN KEY (id_jadwal_ujian)
        REFERENCES jadwal_ujian(id_jadwal_ujian)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_peserta_ujian_siswa
        FOREIGN KEY (id_siswa)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS jawaban_ujian_siswa (
    id_jawaban BIGSERIAL PRIMARY KEY,
    id_peserta_ujian BIGINT NOT NULL,
    id_soal BIGINT NOT NULL,
    id_pilihan BIGINT,
    jawaban_essay TEXT,
    is_benar BOOLEAN DEFAULT NULL,
    waktu_jawab TIMESTAMPTZ,

    CONSTRAINT fk_jawaban_ujian_siswa_peserta_ujian
        FOREIGN KEY (id_peserta_ujian)
        REFERENCES peserta_ujian(id_peserta_ujian)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_jawaban_ujian_siswa_soal
        FOREIGN KEY (id_soal)
        REFERENCES isi_soal(id_soal)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_jawaban_pilihan
        FOREIGN KEY (id_pilihan)
        REFERENCES opsi_pilihan_ganda(id_pilihan_ganda)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT uq_jawaban_peserta_soal
    UNIQUE (id_peserta_ujian, id_soal),

    CONSTRAINT ck_jawaban_xor_nullable
    CHECK (
        (id_pilihan IS NOT NULL AND jawaban_essay IS NULL)
        OR
        (id_pilihan IS NULL AND jawaban_essay IS NOT NULL)
        OR
        (id_pilihan IS NULL AND jawaban_essay IS NULL)
    )
);

-- jadwal_ujian: sering dicari by ujian/sesi/ruangan + waktu
CREATE INDEX IF NOT EXISTS idx_jadwal_ujian_id_ujian
ON jadwal_ujian (id_ujian);

CREATE INDEX IF NOT EXISTS idx_jadwal_ujian_id_sesi
ON jadwal_ujian (id_sesi);

CREATE INDEX IF NOT EXISTS idx_jadwal_ujian_id_ruangan
ON jadwal_ujian (id_ruangan);

CREATE INDEX IF NOT EXISTS idx_jadwal_ujian_waktu_mulai
ON jadwal_ujian (waktu_mulai);


-- peserta_ujian: lookup peserta by jadwal, by siswa, dan status waktu
CREATE INDEX IF NOT EXISTS idx_peserta_ujian_id_jadwal
ON peserta_ujian (id_jadwal_ujian);

CREATE INDEX IF NOT EXISTS idx_peserta_ujian_id_siswa
ON peserta_ujian (id_siswa);

CREATE INDEX IF NOT EXISTS idx_peserta_ujian_waktu_mulai
ON peserta_ujian (waktu_mulai);

CREATE INDEX IF NOT EXISTS idx_peserta_ujian_waktu_submit
ON peserta_ujian (waktu_submit);


-- jawaban_ujian_siswa:
-- UNIQUE(id_peserta_ujian, id_soal) sudah jadi index sendiri, jadi ini opsional:
CREATE INDEX IF NOT EXISTS idx_jawaban_ujian_siswa_id_peserta
ON jawaban_ujian_siswa (id_peserta_ujian);

CREATE INDEX IF NOT EXISTS idx_jawaban_ujian_siswa_id_soal
ON jawaban_ujian_siswa (id_soal);

CREATE INDEX IF NOT EXISTS idx_jawaban_ujian_siswa_id_pilihan
ON jawaban_ujian_siswa (id_pilihan);

CREATE INDEX IF NOT EXISTS idx_jawaban_ujian_siswa_is_benar
ON jawaban_ujian_siswa (is_benar);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS jawaban_ujian_siswa;
DROP TABLE IF EXISTS peserta_ujian;
DROP TABLE IF EXISTS jadwal_ujian;
DROP TABLE IF EXISTS ujian;
-- +goose StatementEnd
