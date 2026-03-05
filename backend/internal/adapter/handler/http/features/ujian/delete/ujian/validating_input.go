package httpx

import coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"

func ValidateInputDeleteUjianRequest(data DeleteUjianRequest) error {
	if data.IDUjian <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
