package bank_soal_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	fakerepo "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/fake_repo"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/update"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBankSoalService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	validID := bank_soal.ID(1)

	tests := []struct {
		name       string
		idBankSoal bank_soal.ID
		patch      updatepatch.UpdateBankSoalPatch
		repo       *fakerepo.FakeBankSoalRepo
		wantErr    error
		wantCalled bool
	}{
		{
			name:       "path 1 -> id bank soal tidak valid",
			idBankSoal: 0,
			patch: updatepatch.UpdateBankSoalPatch{
				NamaBankSoal: ptrString("Nama"),
			},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "path 2 -> no field to update",
			idBankSoal: 1,
			patch:      updatepatch.UpdateBankSoalPatch{},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrNoFieldToUpdate,
			wantCalled: false,
		},
		{
			name:       "path 3 -> id_mapel patch tidak valid",
			idBankSoal: 1,
			patch: updatepatch.UpdateBankSoalPatch{
				IdMapel: ptrID(0),
			},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "path 4 -> id_kelas patch tidak valid",
			idBankSoal: 1,
			patch: updatepatch.UpdateBankSoalPatch{
				IdKelas: ptrID(0),
			},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "path 5 -> id_pengguna patch tidak valid",
			idBankSoal: 1,
			patch: updatepatch.UpdateBankSoalPatch{
				IdPengguna: ptrID(0),
			},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "path 6 -> nama bank soal kosong setelah trim",
			idBankSoal: 1,
			patch: updatepatch.UpdateBankSoalPatch{
				IdPengguna:   &validID,
				NamaBankSoal: ptrString("   "),
			},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCalled: false,
		},
		{
			name:       "path 7 -> deskripsi kosong setelah trim",
			idBankSoal: 1,
			patch: updatepatch.UpdateBankSoalPatch{
				IdPengguna: &validID,
				Deskripsi:  ptrString("   "),
			},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCalled: false,
		},
		{
			name:       "path 8 -> materi kosong setelah trim",
			idBankSoal: 1,
			patch: updatepatch.UpdateBankSoalPatch{
				IdPengguna: &validID,
				Materi:     ptrString("   "),
			},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCalled: false,
		},
		{
			name:       "path 9 -> repository error",
			idBankSoal: 1,
			patch: updatepatch.UpdateBankSoalPatch{
				IdMapel:      ptrID(2),
				IdKelas:      ptrID(3),
				IdPengguna:   &validID,
				NamaBankSoal: ptrString("  UTS  "),
				Deskripsi:    ptrString("  Deskripsi  "),
				Materi:       ptrString("  Materi  "),
			},
			repo:       &fakerepo.FakeBankSoalRepo{UpdateErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       "path 10 -> success",
			idBankSoal: 2,
			patch: updatepatch.UpdateBankSoalPatch{
				IdMapel:      ptrID(2),
				IdKelas:      ptrID(3),
				IdPengguna:   &validID,
				NamaBankSoal: ptrString("  UAS  "),
				Deskripsi:    ptrString("  Deskripsi UAS  "),
				Materi:       ptrString("  Materi UAS  "),
			},
			repo:       &fakerepo.FakeBankSoalRepo{},
			wantErr:    nil,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := bank_soal_service.NewUpdateBankSoalService(tt.repo)
			err := svc.UpdateBankSoalService(ctx, tt.idBankSoal, tt.patch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCalled, tt.repo.UpdateCalled)
		})
	}
}

func ptrString(v string) *string {
	return &v
}

func ptrID(v bank_soal.ID) *bank_soal.ID {
	return &v
}
