-- +goose Up
-- +goose StatementBegin
ALTER TABLE grading_jobs
ADD CONSTRAINT grading_jobs_id_attempt_key UNIQUE (id_attempt);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE grading_jobs
DROP CONSTRAINT IF EXISTS grading_jobs_id_attempt_key;
-- +goose StatementEnd
