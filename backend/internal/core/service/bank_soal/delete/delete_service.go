package bank_soal_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type DeleteBankSoalService struct {
	repo bank_soal_repo.BankSoalRepository
}

func NewDeleteBankSoalService(repo bank_soal_repo.BankSoalRepository) *DeleteBankSoalService {
	return &DeleteBankSoalService{
		repo: repo,
	}
}

func(r *DeleteBankSoalService)DeleteBankSoalService(ctx context.Context, idBankSoal bank_soal.ID) error {
	logger := corelog.FromContext(ctx)

	if err := idChecker(idBankSoal); err != nil {
		logger.Error(ctx, "failed delete bank soal", "layer", "core.service", "op", "bank_soal.delete", "err", coreerror.ErrMissingId)
		return err
	}

	if err := r.repo.DeleteBankSoal(ctx, idBankSoal); err != nil {
		logger.Error(ctx, "failed delete bank soal", "layer", "core.service", "op", "bank_soal.delete", "err", err)
		return err
	}
	return nil
}

// -----------
//  validator
// -----------

func idChecker(Id bank_soal.ID) error {
	if Id <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}