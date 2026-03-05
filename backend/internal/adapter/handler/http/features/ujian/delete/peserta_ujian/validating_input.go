package httpx

import coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"

func ValidateInputDeletePesertaUjianRequest(data DeletePesertaUjianRequest) error {
	if data.IDPesertaUjian <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
