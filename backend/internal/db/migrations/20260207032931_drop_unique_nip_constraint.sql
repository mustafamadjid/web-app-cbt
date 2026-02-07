-- +goose Up
-- +goose StatementBegin
ALTER TABLE profil_guru
DROP CONSTRAINT IF EXISTS uq_guru_nip;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE profil_guru
ADD CONSTRAINT uq_guru_nip UNIQUE (nip);
-- +goose StatementEnd
