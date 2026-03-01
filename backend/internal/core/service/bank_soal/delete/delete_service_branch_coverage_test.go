package bank_soal_service_test

import (
	"context"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/delete"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/delete/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBankSoalService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("branch 1 -> not found from repo diteruskan", func(t *testing.T) {
		repo := &fake_test.FakeBankSoalRepo{DeleteErr: coreerror.ErrNotFound}
		svc := bank_soal_service.NewDeleteBankSoalService(repo)

		err := svc.DeleteBankSoalService(ctx, 1)
		assert.ErrorIs(t, err, coreerror.ErrNotFound)
		assert.True(t, repo.DeleteCalled)
	})

	t.Run("branch 2 -> delete restricted dari repo diteruskan", func(t *testing.T) {
		repo := &fake_test.FakeBankSoalRepo{DeleteErr: coreerror.ErrDeleteRestricted}
		svc := bank_soal_service.NewDeleteBankSoalService(repo)

		err := svc.DeleteBankSoalService(ctx, 2)
		assert.ErrorIs(t, err, coreerror.ErrDeleteRestricted)
		assert.True(t, repo.DeleteCalled)
	})
}
