package bank_soal_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type CreateBankSoalService struct {
	repo bank_soal_repo.BankSoalRepository
}

func NewCreateBankSoalService(repo bank_soal_repo.BankSoalRepository) *CreateBankSoalService {
	return &CreateBankSoalService{repo: repo}
}

func (r *CreateBankSoalService) CreateBankSoalService(ctx context.Context, bankSoal bank_soal.BankSoal) error {
	logger := corelog.FromContext(ctx)

	bankSoal = sanitized(bankSoal)

	if err := idChecker(bankSoal); err != nil {
		logger.Error(ctx, "failed create bank soal", "layer", "core.service", "op", "bank_soal.create", "err", err)
		return err
	}

	if err := r.repo.CreateBankSoal(ctx, bankSoal); err != nil {
		logger.Error(ctx, "failed create bank soal", "layer", "core.service", "op", "bank_soal.create", "err", err)
		return err
	}

	return nil
}

func idChecker(bankSoalId bank_soal.BankSoal) error {
	if bankSoalId.IdMapel <= 0 {
		return coreerror.ErrMissingId
	}

	if bankSoalId.IdKelas <= 0 {
		return coreerror.ErrMissingId
	}

	return nil

}
