-- +goose Up
-- +goose StatementBegin
ALTER TABLE hasil_ujian
ADD COLUMN IF NOT EXISTS essay_graded BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE hasil_ujian
DROP COLUMN IF EXISTS essay_graded;
-- +goose StatementEnd
