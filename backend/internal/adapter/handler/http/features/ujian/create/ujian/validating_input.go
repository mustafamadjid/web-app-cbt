package httpx

import (
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
)

func ValidateInputSafeRequestUjian(data CreatePenjadwalanUjianRequest) error {
	if err := validator.ValidateInputSafe(data.NamaUjian, "nama_ujian"); err != nil {
		return err
	}

	if err := validator.ValidateInputSafe(data.Token, "token"); err != nil {
		
		return err
	}

	if data.DeskripsiUjian != nil {
		if err := validator.ValidateInputSafe(*data.DeskripsiUjian, "deskripsi_ujian"); err != nil {
			return err
		}
	}

	if err := validator.ValidateInputSafe(data.StatusUjian, "status_ujian"); err != nil {
		return err
	}

	return nil
}