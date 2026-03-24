package ujian_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type GradingWorkerRepository interface {
	GetQueuedJobs(ctx context.Context, limit int) ([]ujian.GradingJob, error)
	UpdateStatusJob(ctx context.Context, jobID int, statusJob ujian.JobStatus,errorMsg string, errorCode string) error
}