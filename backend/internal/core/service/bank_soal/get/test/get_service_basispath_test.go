package bank_soal_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	fakerepo "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/fake_repo"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/get"
	"github.com/stretchr/testify/assert"
)

func TestGetBankSoalService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	t.Run("list -> repo error", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{GetBankSoalErr: repoErr}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		_, err := svc.GetBankSoalService(ctx, query.BankSoalFilter{Search: "matematika"})
		assert.ErrorIs(t, err, repoErr)
		assert.True(t, repo.GetCalled)
	})

	t.Run("list -> success", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{
			GetBankSoalData: []bank_soal.BankSoal{
				{IdBankSoal: 1, NamaBankSoal: "UTS", CreatedAt: time.Date(2026, time.March, 2, 10, 0, 0, 0, time.UTC)},
			},
		}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		items, err := svc.GetBankSoalService(ctx, query.BankSoalFilter{
			Search: "  uts  ",
			Limit:  0,
			Offset: -5,
		})
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, "uts", repo.GotFilter.Search)
		assert.Equal(t, 20, repo.GotFilter.Limit)
		assert.Equal(t, 0, repo.GotFilter.Offset)
		assert.Equal(t, "02 Mar 2026", items[0].TanggalDibuat)
	})

	t.Run("list uploaded -> repo error", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{GetBankSoalUploadedErr: repoErr}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		_, err := svc.GetBankSoalUploadedService(ctx, query.BankSoalFilter{Search: "kimia"})
		assert.ErrorIs(t, err, repoErr)
		assert.True(t, repo.GetUploadedCalled)
	})

	t.Run("list uploaded -> success", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{
			GetBankSoalUploadedData: []bank_soal.BankSoal{
				{IdBankSoal: 7, NamaBankSoal: "Uploaded", CreatedAt: time.Date(2026, time.March, 2, 10, 0, 0, 0, time.UTC)},
			},
		}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		items, err := svc.GetBankSoalUploadedService(ctx, query.BankSoalFilter{
			Search: "  uploaded  ",
			Limit:  0,
			Offset: -1,
		})
		assert.NoError(t, err)
		assert.True(t, repo.GetUploadedCalled)
		assert.Equal(t, "uploaded", repo.GotUploadedFilter.Search)
		assert.Equal(t, 20, repo.GotUploadedFilter.Limit)
		assert.Equal(t, 0, repo.GotUploadedFilter.Offset)
		assert.Len(t, items, 1)
		assert.Equal(t, "02 Mar 2026", items[0].TanggalDibuat)
	})

	t.Run("get by id -> invalid id", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		_, err := svc.GetBankSoalByIdService(ctx, 0)
		assert.ErrorIs(t, err, coreerror.ErrMissingId)
		assert.False(t, repo.GetByIDCalled)
	})

	t.Run("get by id -> repo error", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{GetByIDErr: repoErr}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		_, err := svc.GetBankSoalByIdService(ctx, 1)
		assert.ErrorIs(t, err, repoErr)
		assert.True(t, repo.GetByIDCalled)
		assert.Equal(t, bank_soal.ID(1), repo.GotID)
	})

	t.Run("get by id -> not found", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{GetByIDErr: coreerror.ErrNotFound}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		_, err := svc.GetBankSoalByIdService(ctx, 1)
		assert.ErrorIs(t, err, coreerror.ErrNotFound)
		assert.True(t, repo.GetByIDCalled)
	})

	t.Run("get by id -> success", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{
			GetByIDData: bank_soal.BankSoal{
				IdBankSoal:   2,
				NamaBankSoal: "UAS",
				CreatedAt:    time.Date(2026, time.March, 2, 10, 0, 0, 0, time.UTC),
			},
		}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		item, err := svc.GetBankSoalByIdService(ctx, 2)
		assert.NoError(t, err)
		assert.True(t, repo.GetByIDCalled)
		assert.Equal(t, "UAS", item.NamaBankSoal)
		assert.Equal(t, "02 Mar 2026", item.TanggalDibuat)
	})

	t.Run("get by guru -> invalid id", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		_, err := svc.GetBankSoalByGuruService(ctx, 0)
		assert.ErrorIs(t, err, coreerror.ErrMissingId)
		assert.False(t, repo.GetByGuruCalled)
	})

	t.Run("get by guru -> repo error", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{GetByGuruErr: repoErr}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		_, err := svc.GetBankSoalByGuruService(ctx, 9)
		assert.ErrorIs(t, err, repoErr)
		assert.True(t, repo.GetByGuruCalled)
		assert.Equal(t, bank_soal.ID(9), repo.GotGuruID)
	})

	t.Run("get by guru -> success", func(t *testing.T) {
		repo := &fakerepo.FakeBankSoalRepo{
			GetByGuruData: []bank_soal.BankSoal{{
				IdBankSoal:   3,
				NamaBankSoal: "Tryout",
				CreatedAt:    time.Date(2026, time.March, 2, 10, 0, 0, 0, time.UTC),
			}},
		}
		svc := bank_soal_service.NewGetBankSoalService(repo)

		items, err := svc.GetBankSoalByGuruService(ctx, 3)
		assert.NoError(t, err)
		assert.True(t, repo.GetByGuruCalled)
		assert.Len(t, items, 1)
		assert.Equal(t, "02 Mar 2026", items[0].TanggalDibuat)
	})
}
