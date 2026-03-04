package bank_soal_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"strings"
)

func sanitizeNamaBankSoalPatch(payload *updatepatch.UpdateBankSoalPatch) error {
	if payload.NamaBankSoal == nil {
		return nil
	}
	namaBankSoal := strings.TrimSpace(*payload.NamaBankSoal)
	if namaBankSoal == "" {
		return coreerror.ErrMissingField
	}
	payload.NamaBankSoal = &namaBankSoal
	return nil
}
func sanitizeDeskripsiBankSoalPatch(payload *updatepatch.UpdateBankSoalPatch) error {
	if payload.Deskripsi == nil {
		return nil
	}
	deskripsi := strings.TrimSpace(*payload.Deskripsi)
	if deskripsi == "" {
		return coreerror.ErrMissingField
	}
	payload.Deskripsi = &deskripsi
	return nil
}
func sanitizeMateriBankSoalPatch(payload *updatepatch.UpdateBankSoalPatch) error {
	if payload.Materi == nil {
		return nil
	}
	materi := strings.TrimSpace(*payload.Materi)
	if materi == "" {
		return coreerror.ErrMissingField
	}
	payload.Materi = &materi
	return nil
}
