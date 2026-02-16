-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sesi_ujian (
    id_sesi BIGSERIAL PRIMARY KEY,
    kode_sesi VARCHAR(20) NOT NULL,
    nama_sesi VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sesi_ujian_id_sesi ON sesi_ujian(kode_sesi);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sesi_ujian;
-- +goose StatementEnd
