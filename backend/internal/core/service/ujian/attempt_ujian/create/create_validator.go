package attemptujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func validateCreateAttemptUjian(data ujian.AttemptUjian) error {
	if data.IdPesertaUjian <= 0 {
		return coreerror.ErrMissingId
	}

	if data.StatusAttempt == "" {
		return coreerror.ErrMissingField
	}

	if !data.StatusAttempt.ValidStatus() {
		return coreerror.ErrInvalidInput
	}

	if data.WaktuMulai != nil && data.WaktuSubmit != nil && data.WaktuSubmit.Before(*data.WaktuMulai) {
		return coreerror.ErrInvalidInput
	}

	if data.WaktuMulai != nil && data.DeadlineAt != nil && !data.DeadlineAt.After(*data.WaktuMulai) {
		return coreerror.ErrInvalidInput
	}

	if data.StatusAttempt == ujian.ATTEMPT_SUBMITTED && data.WaktuSubmit == nil {
		return coreerror.ErrInvalidInput
	}

	return nil
}
