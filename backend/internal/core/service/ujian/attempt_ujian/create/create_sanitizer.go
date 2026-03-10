package attemptujian_service

import (
	"strings"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func sanitizeCreateAttemptUjian(data ujian.AttemptUjian) ujian.AttemptUjian {
	data.StatusAttempt = ujian.StatusAttempt(strings.ToLower(strings.TrimSpace(string(data.StatusAttempt))))
	if data.StatusAttempt == "" {
		data.StatusAttempt = ujian.ATTEMPT_IN_PROGRESS
	}

	return data
}
