package kelas_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func validateUpdateNamaKelasID(idNamaKelas int) error {
	if idNamaKelas <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateNamaKelasPatch(dataUpdate updatepatch.NamaKelasPatch) error {
	if dataUpdate.IdTingkatKelas == nil && dataUpdate.NamaKelas == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}
func validateIdTingkatKelasPatch(dataUpdate updatepatch.NamaKelasPatch) error {
	if dataUpdate.IdTingkatKelas != nil && *dataUpdate.IdTingkatKelas == 0 {
		return coreerror.ErrMissingField
	}
	return nil
}
