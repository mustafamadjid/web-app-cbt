-- +goose Up
-- +goose StatementBegin
ALTER TABLE statistik_soal
DROP COLUMN IF EXISTS jumlah_soal;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE statistik_soal
ADD COLUMN jumlah_soal INTEGER NOT NULL;
-- +goose StatementEnd
