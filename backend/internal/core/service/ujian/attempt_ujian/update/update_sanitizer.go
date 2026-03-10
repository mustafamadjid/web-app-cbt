package attemptujian_service

import (
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func sanitizeUpdateAttemptStatusPatch(payload *updatepatch.UpdateAttemptUjianPatch) error {
	if payload.StatusAttempt == nil {
		return nil
	}

	status := strings.ToLower(strings.TrimSpace(string(*payload.StatusAttempt)))
	if status == "" {
		return coreerror.ErrMissingField
	}

	normalized := ujian.StatusAttempt(status)
	payload.StatusAttempt = &normalized
	return nil
}
