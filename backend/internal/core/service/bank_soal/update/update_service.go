package bank_soal_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type UpdateBankSoalService struct {
	repo bank_soal_repo.BankSoalRepository
}

func NewUpdateBankSoalService(repo bank_soal_repo.BankSoalRepository) *UpdateBankSoalService {
	return &UpdateBankSoalService{
		repo: repo,
	}
}

func(r *UpdateBankSoalService) UpdateBankSoalService(ctx context.Context, idBankSoal bank_soal.ID, payload bank_soal.BankSoal) (bank_soal.BankSoal, error){
	logger := corelog.FromContext(ctx)

	if idBankSoal <= 0 {
		logger.Error(ctx, "failed update bank soal", "layer", "core.service", "op", "bank_soal.update", "err", coreerror.ErrMissingId)
		return bank_soal.BankSoal{}, nil
	}

	if payload.NamaBankSoal != nil {
		
	}
}