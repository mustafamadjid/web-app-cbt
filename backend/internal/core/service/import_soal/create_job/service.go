package create_job

import (
	"context"
	"fmt"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type CreateJobCmd struct {
	IDBankSoal int64
	IDPengguna int64
	FilePath   string
}

type CreateJobResult struct {
	IDJob int64
}

type CreateJobService struct {
	jobRepo importsoal_repo.ImportSoalJobRepo
}

func NewCreateJobService(jobRepo importsoal_repo.ImportSoalJobRepo) *CreateJobService {
	return &CreateJobService{jobRepo: jobRepo}
}

func (s *CreateJobService) Execute(ctx context.Context, cmd CreateJobCmd) (CreateJobResult, error) {
	logger := corelog.FromContext(ctx)

	if cmd.IDBankSoal <= 0 {
		return CreateJobResult{}, fmt.Errorf("id_bank_soal harus lebih dari 0")
	}

	if cmd.IDPengguna <= 0 {
		return CreateJobResult{}, fmt.Errorf("id_pengguna harus lebih dari 0")
	}

	if cmd.FilePath == "" {
		return CreateJobResult{}, fmt.Errorf("file_path tidak boleh kosong")
	}

	job := importsoal.ImportSoalJob{
		IDBankSoal: cmd.IDBankSoal,
		IDPengguna: cmd.IDPengguna,
		FilePath:   cmd.FilePath,
		Status:     importsoal.StatusPending,
	}

	jobID, err := s.jobRepo.CreateJob(ctx, job)
	if err != nil {
		logger.Error(ctx, "failed creating import job", "layer", "core.service", "op", "import_soal.create_job", "err", err)
		return CreateJobResult{}, err
	}

	return CreateJobResult{IDJob: jobID}, nil
}
