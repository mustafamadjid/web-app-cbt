-- +goose Up
-- +goose StatementBegin
INSERT INTO role (nama_role) VALUES ('ADMIN'), ('GURU'), ('SISWA');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role;
-- +goose StatementEnd
