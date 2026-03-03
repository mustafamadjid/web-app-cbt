package importsoal_repo

import (
	"context"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
)

type ImportSoalJobRepo interface {
	CreateJob(ctx context.Context, job importsoal.ImportSoalJob) (int64, error)
	GetPendingJobs(ctx context.Context, limit int) ([]importsoal.ImportSoalJob, error)
	UpdateJobStatus(ctx context.Context, jobID int64, status importsoal.JobStatus, errorMsg string, totalSoal int) error
	GetJobByID(ctx context.Context, jobID int64) (importsoal.ImportSoalJob, error)
	GetJobsByBankSoal(ctx context.Context, bankSoalID int64) ([]importsoal.ImportSoalJob, error)
}

type ImportBankSoalVersionPayload struct {
	SoalList []importsoal.ParsedSoal
}

type IsiSoalBatchRepo interface {
	ImportBankSoalVersion(ctx context.Context, bankID, userID int64, payload ImportBankSoalVersionPayload) (int64, error)
}
