-- +goose Up
-- +goose StatementBegin
ALTER TABLE jawaban_ujian_siswa
ADD COLUMN IF NOT EXISTS essay_is_benar BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jawaban_ujian_siswa
DROP COLUMN IF EXISTS essay_is_benar;
-- +goose StatementEnd
