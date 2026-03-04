package sesi_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
)

func validateCreateNamaSesi(data sesi.Sesi) error {
	if data.NamaSesi == "" {
		return coreerror.ErrMissingField
	}
	return nil
}
func validateCreateKodeSesi(data sesi.Sesi) error {
	if data.KodeSesi == "" {
		return coreerror.ErrMissingField
	}
	return nil
}
