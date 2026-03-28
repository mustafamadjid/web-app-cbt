package gradingujian_service

import (
	"context"
	"time"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	grading_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/grading"
)

type GradingUjianWorkerRepo struct {
	repo grading_repo.GradingWorkerRepository
	pollInterval time.Duration
}

func NewGradingUjianWorkerService(repo grading_repo.GradingWorkerRepository,pollInterval time.Duration) *GradingUjianWorkerRepo {
	return &GradingUjianWorkerRepo{
		repo: repo,
		pollInterval: pollInterval,
	}
}

func(r *GradingUjianWorkerRepo) Start(ctx context.Context) {
	logger := corelog.FromContext(ctx)

	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()
}

func(r *GradingUjianWorkerRepo) processJobs(ctx context.Context,limit int){
	logger := corelog.FromContext(ctx)
	
	jobs,err := r.repo.GetQueuedJobs(ctx,5)
	if err != nil {
		logger.Error(ctx,"failed fetching pending jobs","layer","worker","err",err)
		return 
	}

	for _,job := range jobs {
		r.processOneJob(ctx,job)
	}

}

func (r *GradingUjianWorkerRepo) processOneJob(ctx context.Context,job ujian.GradingJob){
	logger := corelog.FromContext(ctx)
	

	if err := r.repo.UpdateStatusJob(ctx,job.IDgradingJob,ujian.StatusProcessing,"",""); err != nil {
		logger.Error(ctx,"failed updating status job","layer","worker","err",err)
		return
	}
}
