-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ruang_ujian (
    id_ruangan BIGSERIAL PRIMARY KEY,
    nama_ruangan VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ruang_ujian;
-- +goose StatementEnd
