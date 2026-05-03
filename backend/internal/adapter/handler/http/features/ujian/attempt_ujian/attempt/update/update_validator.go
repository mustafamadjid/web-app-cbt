package httpx

import (
	"errors"
	"strings"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func sanitizeAndValidateUpdateAttemptUjianRequest(req UpdateAttemptUjianRequest) (UpdateAttemptUjianRequest, error) {
	if req.StatusAttempt == nil && req.WaktuSubmit == nil {
		return UpdateAttemptUjianRequest{}, errors.New("no field to update")
	}

	if req.StatusAttempt != nil {
		status := strings.ToLower(strings.TrimSpace(*req.StatusAttempt))
		if status == "" {
			return UpdateAttemptUjianRequest{}, errors.New("status_attempt is required")
		}

		switch ujian.StatusAttempt(status) {
		case ujian.ATTEMPT_SUBMITTED, ujian.ATTEMPT_EXPIRED:
		default:
			return UpdateAttemptUjianRequest{}, errors.New("status_attempt is invalid")
		}

		req.StatusAttempt = &status
	}

	if req.WaktuSubmit != nil && req.WaktuSubmit.IsZero() {
		return UpdateAttemptUjianRequest{}, errors.New("waktu_submit is invalid")
	}

	return req, nil
}
