-- +goose Up
-- +goose StatementBegin
ALTER TABLE jawaban_ujian_siswa
DROP COLUMN IF EXISTS is_benar;

ALTER TABLE jawaban_ujian_siswa
ADD CONSTRAINT unique_attempt_soal
UNIQUE (id_jawaban, id_soal);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jawaban_ujian_siswa
DROP CONSTRAINT IF EXISTS unique_attempt_soal;

ALTER TABLE jawaban_ujian_siswa
ADD COLUMN IF NOT EXISTS is_benar BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd
