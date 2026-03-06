-- +goose Up
-- +goose StatementBegin
ALTER TABLE pengguna
ALTER COLUMN no_hp DROP NOT NULL,
ALTER COLUMN email DROP NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pengguna
ALTER COLUMN no_hp SET NOT NULL,
ALTER COLUMN email SET NOT NULL;
-- +goose StatementEnd
