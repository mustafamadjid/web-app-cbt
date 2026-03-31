-- +goose Up
-- +goose StatementBegin
ALTER TABLE hasil_ujian
DROP CONSTRAINT IF EXISTS graded_id_attempt_and_guru;

ALTER TABLE hasil_ujian
ADD CONSTRAINT graded_id_attempt UNIQUE (id_attempt);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE hasil_ujian
ADD CONSTRAINT graded_id_attempt_and_guru UNIQUE (id_attempt, graded_by);

ALTER TABLE hasil_ujian
DROP CONSTRAINT IF EXISTS graded_id_attempt;
-- +goose StatementEnd
