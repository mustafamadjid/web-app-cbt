package sesi_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"strings"
)

func sanitizeNamaSesiPatch(sesi *updatepatch.UpdateSesiPatch) error {
	if sesi.NamaSesi == nil {
		return nil
	}
	namaSesi := strings.TrimSpace(*sesi.NamaSesi)
	if namaSesi == "" {
		return coreerror.ErrMissingField
	}
	sesi.NamaSesi = &namaSesi
	return nil
}
func sanitizeKodeSesiPatch(sesi *updatepatch.UpdateSesiPatch) error {
	if sesi.KodeSesi == nil {
		return nil
	}
	kodeSesi := strings.TrimSpace(*sesi.KodeSesi)
	if kodeSesi == "" {
		return coreerror.ErrMissingField
	}
	kodeSesi = strings.ToUpper(kodeSesi)
	sesi.KodeSesi = &kodeSesi
	return nil
}
