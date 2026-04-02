package gradingujian_service

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	grading_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/grading"
)

const gradingFailedErrorCode = "GRADING_FAILED"

type GradingUjianExecutor interface {
	GradingUjianPilgan(ctx context.Context, idAttempt int) error
}

type GradingUjianWorkerRepo struct {
	repo         grading_repo.GradingWorkerRepository
	gradingSvc   GradingUjianExecutor
	pollInterval time.Duration
}

func NewGradingUjianWorkerService(
	repo grading_repo.GradingWorkerRepository,
	gradingSvc GradingUjianExecutor,
	pollInterval time.Duration,
) *GradingUjianWorkerRepo {
	return &GradingUjianWorkerRepo{
		repo:         repo,
		gradingSvc:   gradingSvc,
		pollInterval: pollInterval,
	}
}

func (r *GradingUjianWorkerRepo) Start(ctx context.Context) {
	logger := corelog.FromContext(ctx)
	logger.Info(ctx, "grading ujian worker started", "poll_interval", r.pollInterval.String())

	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "grading ujian worker stopped", "err", ctx.Err())
			return
		case <-ticker.C:
			tickCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			func() {
				defer cancel()
				defer func() {
					if recovered := recover(); recovered != nil {
						logger.Error(tickCtx, "panic in grading worker processJobs",
							"panic", recovered,
							"stack", string(debug.Stack()),
						)
					}
				}()

				r.processJobs(tickCtx, 5)
			}()
		}
	}
}

func (r *GradingUjianWorkerRepo) processJobs(ctx context.Context, limit int) {
	logger := corelog.FromContext(ctx)

	jobs, err := r.repo.GetQueuedJobs(ctx, limit)
	if err != nil {
		logger.Error(ctx, "failed fetching pending jobs", "layer", "worker", "err", err)
		return
	}

	for _, job := range jobs {
		r.processOneJob(ctx, job)
	}
}

func (r *GradingUjianWorkerRepo) processOneJob(ctx context.Context, job ujian.GradingJob) {
	logger := corelog.FromContext(ctx).With("grading_job_id", job.IDgradingJob, "attempt_id", job.IDAttempt)

	if err := r.repo.UpdateStatusJob(ctx, job.IDgradingJob, ujian.StatusProcessing, "", ""); err != nil {
		logger.Error(ctx, "failed updating job to processing", "layer", "worker", "err", err)
		return
	}

	if r.gradingSvc == nil {
		err := errors.New("grading ujian service belum terpasang")
		logger.Error(ctx, err.Error(), "layer", "worker")
		if updateErr := r.repo.UpdateStatusJob(ctx, job.IDgradingJob, ujian.StatusFailed, err.Error(), gradingFailedErrorCode); updateErr != nil {
			logger.Error(ctx, "failed updating job to failed", "layer", "worker", "err", updateErr)
		}
		return
	}

	if err := r.gradingSvc.GradingUjianPilgan(ctx, job.IDAttempt); err != nil {
		errMsg := fmt.Sprintf("gagal grading ujian: %v", err)
		logger.Error(ctx, errMsg, "layer", "worker", "err", err)
		if updateErr := r.repo.UpdateStatusJob(ctx, job.IDgradingJob, ujian.StatusFailed, errMsg, gradingFailedErrorCode); updateErr != nil {
			logger.Error(ctx, "failed updating job to failed", "layer", "worker", "err", updateErr)
		}
		return
	}

	if err := r.repo.UpdateStatusJob(ctx, job.IDgradingJob, ujian.StatusDone, "", ""); err != nil {
		logger.Error(ctx, "failed updating job to done", "layer", "worker", "err", err)
		return
	}

	logger.Info(ctx, "grading job completed", "layer", "worker")
}
