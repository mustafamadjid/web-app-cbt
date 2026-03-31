package get_job_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/get_job"
	"github.com/stretchr/testify/assert"
)

func TestGetJobByID(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo error")
	ctx := context.Background()

	tests := []struct {
		name      string
		repo      *FakeImportSoalJobRepo
		jobID     int64
		expectErr error
	}{
		{
			name:      "path 1 -> job id nol atau negatif",
			repo:      &FakeImportSoalJobRepo{},
			jobID:     0,
			expectErr: coreerror.ErrImportJobNotFound,
		},
		{
			name: "path 2 -> gagal get dari repo",
			repo: &FakeImportSoalJobRepo{
				GetJobByIDFn: func(_ context.Context, _ int64) (importsoal.ImportSoalJob, error) {
					return importsoal.ImportSoalJob{}, repoErr
				},
			},
			jobID:     1,
			expectErr: repoErr,
		},
		{
			name: "path 3 -> happy path berhasil get job",
			repo: &FakeImportSoalJobRepo{
				GetJobByIDFn: func(_ context.Context, id int64) (importsoal.ImportSoalJob, error) {
					return importsoal.ImportSoalJob{
						IDJob:      id,
						IDBankSoal: 10,
						Status:     importsoal.StatusCompleted,
						TotalSoal:  5,
						CreatedAt:  time.Now(),
					}, nil
				},
			},
			jobID:     99,
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := get_job.NewGetJobService(tc.repo)
			result, err := svc.GetByID(ctx, tc.jobID)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.jobID, result.IDJob)
			}
		})
	}
}

func TestGetJobsByBankSoal(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo error")
	ctx := context.Background()

	tests := []struct {
		name       string
		repo       *FakeImportSoalJobRepo
		bankSoalID int64
		expectErr  error
		expectLen  int
	}{
		{
			name:       "path 1 -> bank soal id nol atau negatif",
			repo:       &FakeImportSoalJobRepo{},
			bankSoalID: -1,
			expectErr:  coreerror.ErrBankSoalNotFound,
		},
		{
			name: "path 2 -> gagal get dari repo",
			repo: &FakeImportSoalJobRepo{
				GetJobsByBankSoalFn: func(_ context.Context, _ int64) ([]importsoal.ImportSoalJob, error) {
					return nil, repoErr
				},
			},
			bankSoalID: 1,
			expectErr:  repoErr,
		},
		{
			name: "path 3 -> happy path berhasil get jobs",
			repo: &FakeImportSoalJobRepo{
				GetJobsByBankSoalFn: func(_ context.Context, _ int64) ([]importsoal.ImportSoalJob, error) {
					return []importsoal.ImportSoalJob{
						{IDJob: 1, Status: importsoal.StatusCompleted},
						{IDJob: 2, Status: importsoal.StatusPending},
					}, nil
				},
			},
			bankSoalID: 5,
			expectLen:  2,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := get_job.NewGetJobService(tc.repo)
			result, err := svc.GetByBankSoal(ctx, tc.bankSoalID)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tc.expectLen)
			}
		})
	}
}
