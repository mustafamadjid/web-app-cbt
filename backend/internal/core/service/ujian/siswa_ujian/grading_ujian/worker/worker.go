package siswaujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/grading"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"

)

type GradingUjianWorkerRepo struct {
	repo ujian_repo.GradingWorkerRepository
}

func NewGradingUjianWorkerService(repo ujian_repo.GradingWorkerRepository) *GradingUjianWorkerRepo {
	return &GradingUjianWorkerRepo{
		repo: repo,
	}
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