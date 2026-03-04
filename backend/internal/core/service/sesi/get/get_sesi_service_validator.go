package sesi_service

import coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"

func validateSesiID(idSesi int) error {
	if idSesi <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateKodeSesi(kodeSesi string) error {
	if len(kodeSesi) == 0 || kodeSesi == "" {
		return coreerror.ErrMissingField
	}
	return nil
}
