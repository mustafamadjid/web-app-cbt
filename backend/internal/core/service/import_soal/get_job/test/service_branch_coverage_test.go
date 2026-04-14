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

func TestGetJobService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	t.Run("branch 1 -> job id tidak valid", func(t *testing.T) {
		svc := get_job.NewGetJobService(&FakeImportSoalJobRepo{})
		result, err := svc.GetByID(ctx, 0)

		assert.ErrorIs(t, err, coreerror.ErrImportJobNotFound)
		assert.Equal(t, importsoal.ImportSoalJob{}, result)
	})

	t.Run("branch 2 -> repo get job by id gagal", func(t *testing.T) {
		svc := get_job.NewGetJobService(&FakeImportSoalJobRepo{
			GetJobByIDFn: func(context.Context, int64) (importsoal.ImportSoalJob, error) {
				return importsoal.ImportSoalJob{}, repoErr
			},
		})

		result, err := svc.GetByID(ctx, 10)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, importsoal.ImportSoalJob{}, result)
	})

	t.Run("branch 3 -> berhasil get job by id", func(t *testing.T) {
		svc := get_job.NewGetJobService(&FakeImportSoalJobRepo{
			GetJobByIDFn: func(_ context.Context, id int64) (importsoal.ImportSoalJob, error) {
				return importsoal.ImportSoalJob{IDJob: id, CreatedAt: time.Now()}, nil
			},
		})

		result, err := svc.GetByID(ctx, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(10), result.IDJob)
	})

	t.Run("branch 4 -> bank soal id tidak valid", func(t *testing.T) {
		svc := get_job.NewGetJobService(&FakeImportSoalJobRepo{})
		result, err := svc.GetByBankSoal(ctx, 0)

		assert.ErrorIs(t, err, coreerror.ErrBankSoalNotFound)
		assert.Nil(t, result)
	})

	t.Run("branch 5 -> repo get jobs by bank soal gagal", func(t *testing.T) {
		svc := get_job.NewGetJobService(&FakeImportSoalJobRepo{
			GetJobsByBankSoalFn: func(context.Context, int64) ([]importsoal.ImportSoalJob, error) {
				return nil, repoErr
			},
		})

		result, err := svc.GetByBankSoal(ctx, 10)

		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, result)
	})

	t.Run("branch 6 -> berhasil get jobs by bank soal", func(t *testing.T) {
		svc := get_job.NewGetJobService(&FakeImportSoalJobRepo{
			GetJobsByBankSoalFn: func(context.Context, int64) ([]importsoal.ImportSoalJob, error) {
				return []importsoal.ImportSoalJob{{IDJob: 1}, {IDJob: 2}}, nil
			},
		})

		result, err := svc.GetByBankSoal(ctx, 10)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})
}
