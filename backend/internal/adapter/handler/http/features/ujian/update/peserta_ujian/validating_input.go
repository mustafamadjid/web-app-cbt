package httpx

import coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"

func ValidateInputIDRequestUpdatePesertaUjian(data UpdatePesertaUjianRequest) error {
	if data.IdJadwalUjian != nil && *data.IdJadwalUjian <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdSiswa != nil && *data.IdSiswa <= 0 {
		return coreerror.ErrMissingId
	}

	return nil
}
