package matapelajaran_service

import coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"

func validateMapelID(idMapel int) error {
	if idMapel <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
