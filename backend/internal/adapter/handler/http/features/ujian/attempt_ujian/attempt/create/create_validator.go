package httpx

import (
	"errors"
	"strings"

	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
)

func sanitizeAndValidateAttemptUjianRequest(req AttemptUjianRequest) (AttemptUjianRequest, error) {
	if req.IdSiswa <= 0 {
		return AttemptUjianRequest{}, errors.New("id_siswa must be a positive number")
	}
	if req.IdJadwalUjian <= 0 {
		return AttemptUjianRequest{}, errors.New("id_jadwal_ujian must be a positive number")
	}
	if req.WaktuMulai.IsZero() {
		return AttemptUjianRequest{}, errors.New("waktu_mulai is required")
	}

	token, err := validator.ValidateRequiredPrintableText(req.TokenUjian, "token_ujian")
	if err != nil {
		return AttemptUjianRequest{}, err
	}

	req.TokenUjian = strings.ToUpper(token)
	return req, nil
}
