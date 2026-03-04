package ruangujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func validateUpdateRuangUjianID(idRuangan int) error {
	if idRuangan <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateUpdateRuangUjianPatch(ruangUjian updatepatch.UpdateRuangUjianPatch) error {
	if ruangUjian.KodeRuang == nil && ruangUjian.NamaRuang == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}
