package ruangujian_service

import coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"

func validateRuangUjianID(idRuangan int) error {
	if idRuangan <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateKodeRuang(kodeRuang string) error {
	if len(kodeRuang) == 0 || kodeRuang == "" {
		return coreerror.ErrMissingField
	}
	return nil
}
