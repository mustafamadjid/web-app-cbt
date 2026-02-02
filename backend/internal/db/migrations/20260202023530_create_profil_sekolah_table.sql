-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS profil_sekolah (
    id_profil       BIGSERIAL PRIMARY KEY,
    email_sekolah   VARCHAR(50) NOT NULL,
    no_telp_sekolah VARCHAR(20) NOT NULL,
    kepala_sekolah  VARCHAR(100) NOT NULL,
    waka_sekolah    VARCHAR(100) NOT NULL,
    nama_sekolah    VARCHAR(100) NOT NULL,
    alamat_sekolah  VARCHAR(255) NOT NULL,
    logo_sekolah    VARCHAR(255),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS profil_sekolah;
-- +goose StatementEnd
