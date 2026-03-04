package pengumuman_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func validateUpdatePengumumanID(idPengumuman pengumuman.ID) error {
	if idPengumuman <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validatePengumumanUpdateUserID(payload updatepatch.PengumumanUpdatePatch) error {
	if payload.IdPengguna != nil && *payload.IdPengguna <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validatePengumumanUpdatePatch(payload updatepatch.PengumumanUpdatePatch) error {
	if payload.JudulPengumuman == nil && payload.IsiPengumuman == nil && payload.TanggalRilisPengumuman == nil && payload.TanggalSelesaiPengumuman == nil && payload.DokumenPengumuman == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}
