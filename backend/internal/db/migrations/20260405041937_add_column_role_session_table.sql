-- +goose Up
-- +goose StatementBegin

ALTER TABLE sessions
    ADD COLUMN role VARCHAR(100);

UPDATE sessions s
SET role = r.nama_role
FROM pengguna p
JOIN role r ON r.id_role = p.id_role
WHERE p.id_pengguna = s.id_pengguna;

ALTER TABLE sessions
    ALTER COLUMN role SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions
    DROP COLUMN role;
-- +goose StatementEnd
