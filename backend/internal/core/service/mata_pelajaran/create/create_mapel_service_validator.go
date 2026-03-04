package matapelajaran_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
)

func validateMapelCreateInput(mapel matapelajaran.MataPelajaran) error {
	if mapel.IdKelas == 0 {
		return coreerror.ErrInvalidInput
	}
	return nil
}
