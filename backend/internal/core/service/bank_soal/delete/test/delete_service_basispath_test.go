package bank_soal_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/delete"
	fakerepo "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/fake_repo"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBankSoalService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		idBankSoal bank_soal.ID
		repo       *fakerepo.FakeBankSoalRepo
		wantErr    error
		wantCalled bool
	}{
		{
			name:       "path 1 -> id tidak valid",
			idBankSoal: 0,
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "path 2 -> repository error",
			idBankSoal: 10,
			repo:       &fakerepo.FakeBankSoalRepo{DeleteErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       "path 3 -> success",
			idBankSoal: 11,
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    nil,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := bank_soal_service.NewDeleteBankSoalService(tt.repo)
			err := svc.DeleteBankSoalService(ctx, tt.idBankSoal)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCalled, tt.repo.DeleteCalled)
			if tt.wantCalled {
				assert.Equal(t, tt.idBankSoal, tt.repo.GotDeleteID)
			}
		})
	}
}
