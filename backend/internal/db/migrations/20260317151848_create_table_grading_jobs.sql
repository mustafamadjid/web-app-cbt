-- +goose Up
-- +goose StatementBegin
CREATE TABLE grading_jobs (
    id_grading_jobs BIGSERIAL PRIMARY KEY,
    id_attempt BIGINT NOT NULL REFERENCES attempt_ujian(id_attempt) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 3,
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_at TIMESTAMPTZ NULL,
    error_code VARCHAR(100) NULL,
    error_message TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_grading_jobs_status
        CHECK (status IN ('queued', 'processing', 'done', 'failed')),

    CONSTRAINT chk_grading_jobs_retry_count
        CHECK (retry_count >= 0),

    CONSTRAINT chk_grading_jobs_max_retries
        CHECK (max_retries >= 0)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_grading_jobs_fetch
    ON grading_jobs (status, available_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_grading_jobs_attempt_id
    ON grading_jobs (id_attempt);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_grading_jobs_processing_locked_at
    ON grading_jobs (status, locked_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uq_grading_jobs_active_attempt
    ON grading_jobs (id_attempt)
    WHERE status IN ('queued', 'processing');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uq_grading_jobs_active_attempt;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_grading_jobs_processing_locked_at;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_grading_jobs_attempt_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_grading_jobs_fetch;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS grading_jobs;
-- +goose StatementEnd