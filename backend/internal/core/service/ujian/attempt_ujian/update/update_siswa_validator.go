package attemptujian_service

import (
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func validateSiswaUpdateAttemptUjian(idSiswa int, idAttempt ujian.ID, payload updatepatch.UpdateAttemptUjianPatch) error {
	if idSiswa <= 0 || idAttempt <= 0 {
		return coreerror.ErrMissingId
	}

	if payload.IdPesertaUjian != nil || payload.WaktuMulai != nil || payload.DeadlineAt != nil || payload.UpdatedAt != nil {
		return coreerror.ErrInvalidInput
	}

	if payload.StatusAttempt == nil {
		if payload.WaktuSubmit == nil {
			return coreerror.ErrNoFieldToUpdate
		}
		return coreerror.ErrMissingField
	}

	status := ujian.StatusAttempt(strings.ToLower(strings.TrimSpace(string(*payload.StatusAttempt))))
	if status == "" {
		return coreerror.ErrMissingField
	}

	switch status {
	case ujian.ATTEMPT_SUBMITTED:
		if payload.WaktuSubmit == nil {
			return coreerror.ErrMissingField
		}
	case ujian.ATTEMPT_EXPIRED:
		if payload.WaktuSubmit != nil {
			return coreerror.ErrInvalidInput
		}
	default:
		return coreerror.ErrInvalidInput
	}

	return nil
}
