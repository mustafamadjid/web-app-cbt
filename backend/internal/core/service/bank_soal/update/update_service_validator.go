package bank_soal_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func validateUpdateBankSoalID(idBankSoal bank_soal.ID) error {
	if idBankSoal <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateUpdateBankSoalPatch(payload updatepatch.UpdateBankSoalPatch) error {
	if payload.IdMapel == nil && payload.IdKelas == nil && payload.IdPengguna == nil && payload.NamaBankSoal == nil && payload.Deskripsi == nil && payload.Materi == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}
func validateUpdateBankSoalPatchID(payload updatepatch.UpdateBankSoalPatch) error {
	if payload.IdKelas != nil && *payload.IdKelas <= 0 {
		return coreerror.ErrMissingId
	}
	if payload.IdPengguna != nil && *payload.IdPengguna <= 0 {
		return coreerror.ErrMissingId
	}
	if payload.IdMapel != nil && *payload.IdMapel <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
