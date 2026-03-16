package httpx

import (
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

func ValidateUpdateUjianRequestFields(data UpdatePenjadwalanUjianRequest) (UpdatePenjadwalanUjianRequest, error) {
	if data.NamaUjian != nil {
		value, err := validator.ValidateRequiredPrintableText(*data.NamaUjian, "nama_ujian")
		if err != nil {
			return UpdatePenjadwalanUjianRequest{}, err
		}
		data.NamaUjian = &value
	}

	if data.Token != nil {
		value, err := validator.ValidateRequiredPrintableText(*data.Token, "token")
		if err != nil {
			return UpdatePenjadwalanUjianRequest{}, err
		}
		data.Token = &value
	}

	if data.DeskripsiUjian != nil {
		value, err := validator.ValidateRequiredPrintableText(*data.DeskripsiUjian, "deskripsi_ujian")
		if err != nil {
			return UpdatePenjadwalanUjianRequest{}, err
		}
		data.DeskripsiUjian = &value
	}

	if data.StatusUjian != nil {
		value, err := validator.ValidateRequiredPrintableText(*data.StatusUjian, "status_ujian")
		if err != nil {
			return UpdatePenjadwalanUjianRequest{}, err
		}
		data.StatusUjian = &value
	}

	return data, nil
}

func ValidateUpdateUjianRequestIDs(data UpdatePenjadwalanUjianRequest) error {
	if data.IdBankSoal != nil && *data.IdBankSoal <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdKelas != nil && *data.IdKelas <= 0 {
		return coreerror.ErrMissingId
	}

	if data.IdNamaKelas != nil && *data.IdNamaKelas < 0 {
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
