package create_job_test

import (
	"context"
	"errors"
	"testing"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/create_job"
	"github.com/stretchr/testify/assert"
)

func TestCreateJobService_BranchCoverage(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo error")
	ctx := context.Background()

	tests := []struct {
		name       string
		repo       *FakeImportSoalJobRepo
		cmd        create_job.CreateJobCmd
		wantErr    bool
		wantID     int64
		wantCalled bool
	}{
		{
			name:       "branch 1 -> id bank soal tidak valid",
			repo:       &FakeImportSoalJobRepo{},
			cmd:        create_job.CreateJobCmd{IDBankSoal: 0, IDPengguna: 1, FilePath: "/tmp/test.docx"},
			wantErr:    true,
			wantCalled: false,
		},
		{
			name:       "branch 2 -> id pengguna tidak valid",
			repo:       &FakeImportSoalJobRepo{},
			cmd:        create_job.CreateJobCmd{IDBankSoal: 1, IDPengguna: 0, FilePath: "/tmp/test.docx"},
			wantErr:    true,
			wantCalled: false,
		},
		{
			name:       "branch 3 -> file path kosong",
			repo:       &FakeImportSoalJobRepo{},
			cmd:        create_job.CreateJobCmd{IDBankSoal: 1, IDPengguna: 1},
			wantErr:    true,
			wantCalled: false,
		},
		{
			name: "branch 4 -> repo create job gagal",
			repo: &FakeImportSoalJobRepo{
				CreateJobFn: func(_ context.Context, _ importsoal.ImportSoalJob) (int64, error) {
					return 0, repoErr
				},
			},
			cmd:        create_job.CreateJobCmd{IDBankSoal: 1, IDPengguna: 2, FilePath: "/tmp/test.docx"},
			wantErr:    true,
			wantCalled: true,
		},
		{
			name: "branch 5 -> berhasil create job",
			repo: &FakeImportSoalJobRepo{
				CreateJobFn: func(_ context.Context, job importsoal.ImportSoalJob) (int64, error) {
					assert.Equal(t, int64(1), job.IDBankSoal)
					assert.Equal(t, int64(2), job.IDPengguna)
					assert.Equal(t, importsoal.StatusPending, job.Status)
					return 44, nil
				},
			},
			cmd:        create_job.CreateJobCmd{IDBankSoal: 1, IDPengguna: 2, FilePath: "/tmp/test.docx"},
			wantID:     44,
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := create_job.NewCreateJobService(tc.repo)
			result, err := svc.Execute(ctx, tc.cmd)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantID, result.IDJob)
			assert.Equal(t, tc.wantCalled, tc.repo.CreateJobCalled)
		})
	}
}
