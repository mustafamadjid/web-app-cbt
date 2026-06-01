package bank_soal_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/create"
	fakerepo "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/fake_repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBankSoalService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	validInput := bank_soal.BankSoal{
		IdMapel:      10,
		IdKelas:      11,
		IdPengguna:   12,
		NamaBankSoal: "  Nama Bank Soal  ",
		Deskripsi:    "  Deskripsi  ",
		Materi:       "  Materi  ",
	}

	tests := []struct {
		name       string
		repo       *fakerepo.FakeBankSoalRepo
		input      bank_soal.BankSoal
		wantErr    error
		wantCalled bool
		wantCreate bank_soal.BankSoal
	}{
		{
			name:       "Branch 1 -> IdMapel sama dengan nol mengembalikan ErrMissingId dan tidak memanggil repository",
			repo:       &fakerepo.FakeBankSoalRepo{},
			input:      bank_soal.BankSoal{IdMapel: 0, IdKelas: 11, IdPengguna: 12},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "Branch 2 -> IdMapel negatif mengembalikan ErrMissingId dan tidak memanggil repository",
			repo:       &fakerepo.FakeBankSoalRepo{},
			input:      bank_soal.BankSoal{IdMapel: -1, IdKelas: 11, IdPengguna: 12},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "Branch 3 -> IdKelas sama dengan nol mengembalikan ErrMissingId dan tidak memanggil repository",
			repo:       &fakerepo.FakeBankSoalRepo{},
			input:      bank_soal.BankSoal{IdMapel: 10, IdKelas: 0, IdPengguna: 12},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "Branch 4 -> IdKelas negatif mengembalikan ErrMissingId dan tidak memanggil repository",
			repo:       &fakerepo.FakeBankSoalRepo{},
			input:      bank_soal.BankSoal{IdMapel: 10, IdKelas: -1, IdPengguna: 12},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "Branch 5 -> repository gagal mengembalikan error repository",
			repo:       &fakerepo.FakeBankSoalRepo{CreateErr: repoErr},
			input:      validInput,
			wantErr:    repoErr,
			wantCalled: true,
			wantCreate: bank_soal.BankSoal{
				IdMapel:      10,
				IdKelas:      11,
				IdPengguna:   12,
				NamaBankSoal: "Nama Bank Soal",
				Deskripsi:    "Deskripsi",
				Materi:       "Materi",
			},
		},
		{
			name:       "Branch 6 -> semua input valid berhasil membuat bank soal",
			repo:       &fakerepo.FakeBankSoalRepo{},
			input:      validInput,
			wantErr:    nil,
			wantCalled: true,
			wantCreate: bank_soal.BankSoal{
				IdMapel:      10,
				IdKelas:      11,
				IdPengguna:   12,
				NamaBankSoal: "Nama Bank Soal",
				Deskripsi:    "Deskripsi",
				Materi:       "Materi",
			},
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
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantCalled, tt.repo.CreateCalled)
			if tt.wantCalled {
				assert.Equal(t, tt.wantCreate, tt.repo.GotCreate)
			}
		})
	}
}
