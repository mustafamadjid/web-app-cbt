package bank_soal_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
)

func validateBankSoalID(id bank_soal.ID) error {
	if id <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
