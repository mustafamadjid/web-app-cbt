package bank_soal_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/create"
	"github.com/stretchr/testify/assert"
)

func TestCreateBankSoalService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		repo       *FakeBankSoalRepo
		input      bank_soal.BankSoal
		wantErr    error
		wantCalled bool
	}{
		{
			name: "path 1 -> id mapel tidak valid",
			repo: &FakeBankSoalRepo{},
			input: bank_soal.BankSoal{
				IdMapel:      0,
				IdKelas:      1,
				NamaBankSoal: "Bank Soal",
			},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: "path 2 -> id kelas tidak valid",
			repo: &FakeBankSoalRepo{},
			input: bank_soal.BankSoal{
				IdMapel:      1,
				IdKelas:      0,
				NamaBankSoal: "Bank Soal",
			},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: "path 3 -> repository error",
			repo: &FakeBankSoalRepo{CreateErr: repoErr},
			input: bank_soal.BankSoal{
				IdMapel:      2,
				IdKelas:      3,
				IdPengguna:   4,
				NamaBankSoal: "  Bank Soal UTS  ",
				Deskripsi:    "  Deskripsi UTS  ",
				Materi:       "  Bab 1  ",
			},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name: "path 4 -> berhasil create",
			repo: &FakeBankSoalRepo{},
			input: bank_soal.BankSoal{
				IdMapel:      2,
				IdKelas:      3,
				IdPengguna:   4,
				NamaBankSoal: "  Bank Soal UAS  ",
				Deskripsi:    "  Deskripsi UAS  ",
				Materi:       "  Bab 2  ",
			},
			wantErr:    nil,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := bank_soal_service.NewCreateBankSoalService(tt.repo)
			err := svc.CreateBankSoalService(ctx, tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCalled, tt.repo.CreateCalled)
		})
	}
}
