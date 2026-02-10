-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS aktivitas_user_created_at_desc_idx
ON aktivitas_user(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS aktivitas_user_created_at_desc_idx;
-- +goose StatementEnd
