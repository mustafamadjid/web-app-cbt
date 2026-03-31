package create_job_test

import (
	"context"
	"errors"
	"testing"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/create_job"
	"github.com/stretchr/testify/assert"
)

func TestCreateJobService(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo error")
	ctx := context.Background()

	tests := []struct {
		name      string
		repo      *FakeImportSoalJobRepo
		cmd       create_job.CreateJobCmd
		expectErr bool
		expectID  int64
	}{
		{
			name: "path 1 -> id_bank_soal kosong atau nol",
			repo: &FakeImportSoalJobRepo{},
			cmd: create_job.CreateJobCmd{
				IDBankSoal: 0,
				IDPengguna: 1,
				FilePath:   "/tmp/test.docx",
			},
			expectErr: true,
		},
		{
			name: "path 2 -> id_pengguna kosong atau nol",
			repo: &FakeImportSoalJobRepo{},
			cmd: create_job.CreateJobCmd{
				IDBankSoal: 1,
				IDPengguna: 0,
				FilePath:   "/tmp/test.docx",
			},
			expectErr: true,
		},
		{
			name: "path 3 -> file_path kosong",
			repo: &FakeImportSoalJobRepo{},
			cmd: create_job.CreateJobCmd{
				IDBankSoal: 1,
				IDPengguna: 1,
				FilePath:   "",
			},
			expectErr: true,
		},
		{
			name: "path 4 -> gagal create job di repo",
			repo: &FakeImportSoalJobRepo{
				CreateJobFn: func(_ context.Context, _ importsoal.ImportSoalJob) (int64, error) {
					return 0, repoErr
				},
			},
			cmd: create_job.CreateJobCmd{
				IDBankSoal: 1,
				IDPengguna: 1,
				FilePath:   "/tmp/test.docx",
			},
			expectErr: true,
		},
		{
			name: "path 5 -> happy path berhasil create job",
			repo: &FakeImportSoalJobRepo{
				CreateJobFn: func(_ context.Context, job importsoal.ImportSoalJob) (int64, error) {
					assert.Equal(t, int64(5), job.IDBankSoal)
					assert.Equal(t, int64(2), job.IDPengguna)
					assert.Equal(t, importsoal.StatusPending, job.Status)
					return 42, nil
				},
			},
			cmd: create_job.CreateJobCmd{
				IDBankSoal: 5,
				IDPengguna: 2,
				FilePath:   "/tmp/test.docx",
			},
			expectErr: false,
			expectID:  42,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := create_job.NewCreateJobService(tc.repo)
			result, err := svc.Execute(ctx, tc.cmd)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectID, result.IDJob)
			}
		})
	}
}
