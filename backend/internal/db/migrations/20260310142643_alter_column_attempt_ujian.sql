-- +goose Up
-- +goose StatementBegin
ALTER TABLE attempt_ujian
    DROP COLUMN attempt_no;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE attempt_ujian
    ADD COLUMN attempt_no INTEGER NOT NULL;
-- +goose StatementEnd
