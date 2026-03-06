-- +goose Up
-- +goose StatementBegin
ALTER TABLE pengguna
DROP CONSTRAINT uq_pengguna_email;

ALTER TABLE pengguna
DROP CONSTRAINT uq_pengguna_no_hp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pengguna
ADD CONSTRAINT uq_pengguna_email UNIQUE (email);

ALTER TABLE pengguna
ADD CONSTRAINT uq_pengguna_no_hp UNIQUE (no_hp);
-- +goose StatementEnd
