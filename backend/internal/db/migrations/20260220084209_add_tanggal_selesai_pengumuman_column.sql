-- +goose Up
-- +goose StatementBegin
ALTER TABLE pengumuman
    ADD COLUMN tanggal_selesai_pengumuman DATE NOT NULL,
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pengumuman
    DROP COLUMN IF EXISTS tanggal_selesai_pengumuman,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
