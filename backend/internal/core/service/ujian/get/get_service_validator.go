package ujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func validateUjianID(id ujian.ID) error {
	if id <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validatePesertaFilter(filter ujian.PesertaUjian) error {
	if filter.IdPesertaUjian < 0 || filter.IdJadwalUjian < 0 || filter.IdSiswa < 0 {
		return errInvalidPeserta
	}
	return nil
}
