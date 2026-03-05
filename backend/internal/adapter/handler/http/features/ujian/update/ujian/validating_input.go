package httpx

import (
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

func ValidateInputSafeRequestUpdateUjian(data UpdatePenjadwalanUjianRequest) error {
	if data.NamaUjian != nil {
		if err := validator.ValidateInputSafe(*data.NamaUjian, "nama_ujian"); err != nil {
			return err
		}
	}

	if data.Token != nil {
		if err := validator.ValidateInputSafe(*data.Token, "token"); err != nil {
			return err
		}
	}

	if data.DeskripsiUjian != nil {
		if err := validator.ValidateInputSafe(*data.DeskripsiUjian, "deskripsi_ujian"); err != nil {
			return err
		}
	}

	if data.StatusUjian != nil {
		if err := validator.ValidateInputSafe(*data.StatusUjian, "status_ujian"); err != nil {
			return err
		}
	}

	return nil
}

func ValidateInputIDRequestUpdateUjian(data UpdatePenjadwalanUjianRequest) error {
	if data.IdBankSoal != nil && *data.IdBankSoal <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdKelas != nil && *data.IdKelas <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdNamaKelas != nil && *data.IdNamaKelas <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdGuru != nil && *data.IdGuru <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdSesi != nil && *data.IdSesi <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdRuangan != nil && *data.IdRuangan <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdPengawas != nil && *data.IdPengawas <= 0 {
		return coreerror.ErrMissingId
	}

	return nil
}
