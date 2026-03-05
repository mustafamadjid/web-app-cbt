package httpx

import coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"

func ValidateInputIdRequestPesertUjian(data CreatePesertaUjianRequest) error {
	if data.IdJadwalUjian <= 0 || data.IdSiswa <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}