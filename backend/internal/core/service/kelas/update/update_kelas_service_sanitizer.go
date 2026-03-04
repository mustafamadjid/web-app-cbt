package kelas_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"strings"
)

func sanitizeNamaKelasPatch(dataUpdate *updatepatch.NamaKelasPatch) error {
	if dataUpdate.NamaKelas != nil {
		trimmedNamaKelas := strings.TrimSpace(*dataUpdate.NamaKelas)
		if trimmedNamaKelas == "" {
			return coreerror.ErrMissingField
		}
		dataUpdate.NamaKelas = &trimmedNamaKelas
	}
	return nil
}
