-- +goose Up
-- +goose StatementBegin
ALTER TABLE peserta_ujian
DROP COLUMN waktu_mulai,
DROP COLUMN waktu_submit,
DROP COLUMN nilai_ujian;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE peserta_ujian
ADD COLUMN waktu_mulai TIMESTAMPTZ,
ADD COLUMN waktu_submit TIMESTAMPTZ,
ADD COLUMN nilai_ujian NUMERIC(5,2);
-- +goose StatementEnd
