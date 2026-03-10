package attemptujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func validateUpdateAttemptUjian(idAttempt ujian.ID, payload updatepatch.UpdateAttemptUjianPatch) error {
	if idAttempt <= 0 {
		return coreerror.ErrMissingId
	}

	if payload.IdPesertaUjian == nil &&
		payload.StatusAttempt == nil &&
		payload.WaktuMulai == nil &&
		payload.WaktuSubmit == nil &&
		payload.DeadlineAt == nil {
		return coreerror.ErrNoFieldToUpdate
	}

	if payload.IdPesertaUjian != nil && *payload.IdPesertaUjian <= 0 {
		return coreerror.ErrMissingId
	}

	if payload.StatusAttempt != nil && !payload.StatusAttempt.ValidStatus() {
		return coreerror.ErrInvalidInput
	}

	if payload.WaktuMulai != nil && payload.WaktuMulai.IsZero() {
		return coreerror.ErrInvalidInput
	}

	if payload.WaktuSubmit != nil && payload.WaktuSubmit.IsZero() {
		return coreerror.ErrInvalidInput
	}

	if payload.DeadlineAt != nil && payload.DeadlineAt.IsZero() {
		return coreerror.ErrInvalidInput
	}

	if payload.WaktuMulai != nil && payload.WaktuSubmit != nil && payload.WaktuSubmit.Before(*payload.WaktuMulai) {
		return coreerror.ErrInvalidInput
	}

	if payload.WaktuMulai != nil && payload.DeadlineAt != nil && !payload.DeadlineAt.After(*payload.WaktuMulai) {
		return coreerror.ErrInvalidInput
	}

	return nil
}
