package httpx

import (
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
)

func ValidateCreateUjianRequestFields(data CreatePenjadwalanUjianRequest) (CreatePenjadwalanUjianRequest, error) {
	namaUjian, err := validator.ValidateRequiredPrintableText(data.NamaUjian, "nama_ujian")
	if err != nil {
		return CreatePenjadwalanUjianRequest{}, err
	}
	data.NamaUjian = namaUjian

	token, err := validator.ValidateRequiredPrintableText(data.Token, "token")
	if err != nil {
		return CreatePenjadwalanUjianRequest{}, err
	}
	data.Token = token

	if data.DeskripsiUjian != nil {
		value, err := validator.ValidateRequiredPrintableText(*data.DeskripsiUjian, "deskripsi_ujian")
		if err != nil {
			return CreatePenjadwalanUjianRequest{}, err
		}
		data.DeskripsiUjian = &value
	}

	statusUjian, err := validator.ValidateRequiredPrintableText(data.StatusUjian, "status_ujian")
	if err != nil {
		return CreatePenjadwalanUjianRequest{}, err
	}
	data.StatusUjian = statusUjian

	return data, nil
}
