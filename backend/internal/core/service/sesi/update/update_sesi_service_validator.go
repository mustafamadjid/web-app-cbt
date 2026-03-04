package sesi_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func validateUpdateSesiID(idSesi int) error {
	if idSesi == 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateUpdateSesiPatch(sesi updatepatch.UpdateSesiPatch) error {
	if sesi.KodeSesi == nil && sesi.NamaSesi == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}
