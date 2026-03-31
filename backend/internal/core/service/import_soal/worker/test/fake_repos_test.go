package worker_test

import (
	"context"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoalrepo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
)

// FakeImportSoalJobRepo is a fake implementation of importsoal_repo.ImportSoalJobRepo
type FakeImportSoalJobRepo struct {
	CreateJobFn         func(ctx context.Context, job importsoal.ImportSoalJob) (int64, error)
	GetPendingJobsFn    func(ctx context.Context, limit int) ([]importsoal.ImportSoalJob, error)
	UpdateJobStatusFn   func(ctx context.Context, jobID int64, status importsoal.JobStatus, errorMsg string, totalSoal int) error
	GetJobByIDFn        func(ctx context.Context, jobID int64) (importsoal.ImportSoalJob, error)
	GetJobsByBankSoalFn func(ctx context.Context, bankSoalID int64) ([]importsoal.ImportSoalJob, error)
}

func (f *FakeImportSoalJobRepo) CreateJob(ctx context.Context, job importsoal.ImportSoalJob) (int64, error) {
	if f.CreateJobFn != nil {
		return f.CreateJobFn(ctx, job)
	}
	return 1, nil
}

func (f *FakeImportSoalJobRepo) GetPendingJobs(ctx context.Context, limit int) ([]importsoal.ImportSoalJob, error) {
	if f.GetPendingJobsFn != nil {
		return f.GetPendingJobsFn(ctx, limit)
	}
	return nil, nil
}

func (f *FakeImportSoalJobRepo) UpdateJobStatus(ctx context.Context, jobID int64, status importsoal.JobStatus, errorMsg string, totalSoal int) error {
	if f.UpdateJobStatusFn != nil {
		return f.UpdateJobStatusFn(ctx, jobID, status, errorMsg, totalSoal)
	}
	return nil
}

func (f *FakeImportSoalJobRepo) GetJobByID(ctx context.Context, jobID int64) (importsoal.ImportSoalJob, error) {
	if f.GetJobByIDFn != nil {
		return f.GetJobByIDFn(ctx, jobID)
	}
	return importsoal.ImportSoalJob{}, nil
}

func (f *FakeImportSoalJobRepo) GetJobsByBankSoal(ctx context.Context, bankSoalID int64) ([]importsoal.ImportSoalJob, error) {
	if f.GetJobsByBankSoalFn != nil {
		return f.GetJobsByBankSoalFn(ctx, bankSoalID)
	}
	return nil, nil
}

// FakeIsiSoalBatchRepo is a fake implementation of importsoal_repo.IsiSoalBatchRepo
type FakeIsiSoalBatchRepo struct {
	ImportBankSoalVersionFn func(ctx context.Context, bankID, userID int64, payload importsoalrepo.ImportBankSoalVersionPayload) (int64, error)
}

func (f *FakeIsiSoalBatchRepo) ImportBankSoalVersion(ctx context.Context, bankID, userID int64, payload importsoalrepo.ImportBankSoalVersionPayload) (int64, error) {
	if f.ImportBankSoalVersionFn != nil {
		return f.ImportBankSoalVersionFn(ctx, bankID, userID, payload)
	}
	return 1, nil
}
