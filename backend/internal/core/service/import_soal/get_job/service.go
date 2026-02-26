package get_job

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type GetJobService struct {
	jobRepo importsoal_repo.ImportSoalJobRepo
}

func NewGetJobService(jobRepo importsoal_repo.ImportSoalJobRepo) *GetJobService {
	return &GetJobService{jobRepo: jobRepo}
}

func (s *GetJobService) GetByID(ctx context.Context, jobID int64) (importsoal.ImportSoalJob, error) {
	logger := corelog.FromContext(ctx)

	if jobID <= 0 {
		return importsoal.ImportSoalJob{}, coreerror.ErrImportJobNotFound
	}

	job, err := s.jobRepo.GetJobByID(ctx, jobID)
	if err != nil {
		logger.Error(ctx, "failed getting import job", "layer", "core.service", "op", "import_soal.get_job", "err", err)
		return importsoal.ImportSoalJob{}, err
	}

	return job, nil
}

func (s *GetJobService) GetByBankSoal(ctx context.Context, bankSoalID int64) ([]importsoal.ImportSoalJob, error) {
	logger := corelog.FromContext(ctx)

	if bankSoalID <= 0 {
		return nil, coreerror.ErrBankSoalNotFound
	}

	jobs, err := s.jobRepo.GetJobsByBankSoal(ctx, bankSoalID)
	if err != nil {
		logger.Error(ctx, "failed getting import jobs by bank soal", "layer", "core.service", "op", "import_soal.get_jobs_by_bank_soal", "err", err)
		return nil, err
	}

	return jobs, nil
}
