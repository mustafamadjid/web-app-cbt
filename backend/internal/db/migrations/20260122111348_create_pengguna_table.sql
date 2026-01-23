-- +goose Up
-- +goose StatementBegin

-- 1) Tabel role
CREATE TABLE IF NOT EXISTS role (
    id_role     BIGSERIAL PRIMARY KEY,
    nama_role   VARCHAR(50) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_role_nama UNIQUE (nama_role)
);

-- 2) Tabel kelas
CREATE TABLE IF NOT EXISTS kelas (
    id_kelas       BIGSERIAL PRIMARY KEY,
    tingkat_kelas  SMALLINT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT ck_kelas_tingkat CHECK (tingkat_kelas > 0),
    CONSTRAINT uq_kelas_tingkat UNIQUE (tingkat_kelas)
);

-- 2b) Tabel nama_kelas 
CREATE TABLE IF NOT EXISTS nama_kelas (
    id_nama_kelas  BIGSERIAL PRIMARY KEY,
    id_kelas       BIGINT NOT NULL,
    nama_kelas     VARCHAR(100) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_nama_kelas_kelas
        FOREIGN KEY (id_kelas)
        REFERENCES kelas(id_kelas)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT uq_nama_kelas_per_tingkat UNIQUE (id_kelas, nama_kelas)
);

-- 3) Tabel pengguna
CREATE TABLE IF NOT EXISTS pengguna (
    id_pengguna   BIGSERIAL PRIMARY KEY,
    foto          VARCHAR(255),
    nama_lengkap  VARCHAR(150) NOT NULL,
    jenis_kelamin SMALLINT NOT NULL, -- 1=LK, 2=PR
    username      VARCHAR(50)  NOT NULL,
    password      VARCHAR(255) NOT NULL, -- hash password
    email         VARCHAR(150) NOT NULL,
    no_hp         VARCHAR(20),
    id_role       BIGINT NOT NULL,
    status_akun   VARCHAR(30)  NOT NULL, -- AKTIF | NONAKTIF
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_pengguna_username UNIQUE (username),
    CONSTRAINT uq_pengguna_email UNIQUE (email),
    CONSTRAINT uq_pengguna_no_hp UNIQUE (no_hp),
    CONSTRAINT ck_pengguna_jenis_kelamin CHECK (jenis_kelamin IN (1, 2)),

    CONSTRAINT fk_pengguna_role
        FOREIGN KEY (id_role)
        REFERENCES role(id_role)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

-- 4) Tabel profil_siswa
CREATE TABLE IF NOT EXISTS profil_siswa (
    id_siswa      BIGSERIAL PRIMARY KEY,
    id_pengguna   BIGINT NOT NULL,
    id_kelas      BIGINT NOT NULL,
    nisn          VARCHAR(20) NOT NULL,
    no_absen      INT NOT NULL,
    angkatan      INT NOT NULL,
    tempat_lahir  VARCHAR(255) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_siswa_nisn UNIQUE (nisn),
    CONSTRAINT uq_siswa_pengguna UNIQUE (id_pengguna), -- 1 user = 1 profil siswa

    CONSTRAINT fk_siswa_pengguna
        FOREIGN KEY (id_pengguna)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_siswa_kelas
        FOREIGN KEY (id_kelas)
        REFERENCES kelas(id_kelas)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT ck_siswa_no_absen CHECK (no_absen > 0),
    CONSTRAINT ck_siswa_angkatan CHECK (angkatan > 0)
);

-- 5) Tabel profil_guru
CREATE TABLE IF NOT EXISTS profil_guru (
    id_guru      BIGSERIAL PRIMARY KEY,
    id_pengguna  BIGINT NOT NULL,
    nip          VARCHAR(20) NOT NULL,
    jabatan      VARCHAR(100) NOT NULL,
    bidang_studi VARCHAR(100) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_guru_nip UNIQUE (nip),
    CONSTRAINT uq_guru_pengguna UNIQUE (id_pengguna), -- 1 user = 1 profil guru

    CONSTRAINT fk_guru_pengguna
        FOREIGN KEY (id_pengguna)
        REFERENCES pengguna(id_pengguna)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

-- =========================
-- Indexing
-- =========================

-- nama_kelas
CREATE INDEX IF NOT EXISTS idx_nama_kelas_id_kelas ON nama_kelas(id_kelas);

-- pengguna
CREATE INDEX IF NOT EXISTS idx_pengguna_id_role ON pengguna(id_role);
CREATE INDEX IF NOT EXISTS idx_pengguna_status_akun ON pengguna(status_akun);


-- profil_siswa
CREATE INDEX IF NOT EXISTS idx_siswa_id_kelas ON profil_siswa(id_kelas);
CREATE INDEX IF NOT EXISTS idx_siswa_angkatan ON profil_siswa(angkatan);
CREATE INDEX IF NOT EXISTS idx_siswa_kelas_absen ON profil_siswa(id_kelas, no_absen);

-- profil_guru
CREATE INDEX IF NOT EXISTS idx_guru_bidang_studi ON profil_guru(bidang_studi);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS profil_guru;
DROP TABLE IF EXISTS profil_siswa;
DROP TABLE IF EXISTS pengguna;
DROP TABLE IF EXISTS nama_kelas;
DROP TABLE IF EXISTS kelas;
DROP TABLE IF EXISTS role;

-- +goose StatementEnd
