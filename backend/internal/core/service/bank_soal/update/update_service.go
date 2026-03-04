package bank_soal_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateBankSoalService struct {
	repo bank_soal_repo.BankSoalRepository
}

func NewUpdateBankSoalService(repo bank_soal_repo.BankSoalRepository) *UpdateBankSoalService {
	return &UpdateBankSoalService{
		repo: repo,
	}
}

func (r *UpdateBankSoalService) UpdateBankSoalService(ctx context.Context, idBankSoal bank_soal.ID, payload updatepatch.UpdateBankSoalPatch) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdateBankSoalID(idBankSoal); err != nil {
		logger.Error(ctx, "failed update bank soal", "layer", "core.service", "op", "bank_soal.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateBankSoalPatch(payload); err != nil {
		logger.Error(ctx, "failed update bank soal", "layer", "core.service", "op", "bank_soal.update", "err", err)
		return err
	}

	if err := validateUpdateBankSoalPatchID(payload); err != nil {
		logger.Error(ctx, "failed update bank soal", "layer", "core.service", "op", "bank_soal.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := sanitizeNamaBankSoalPatch(&payload); err != nil {
		logger.Error(ctx, "failed update bank soal", "layer", "core.service", "op", "bank_soal.update", "err", err)
		return err
	}

	if err := sanitizeDeskripsiBankSoalPatch(&payload); err != nil {
		logger.Error(ctx, "failed update bank soal", "layer", "core.service", "op", "bank_soal.update", "err", err)
		return err
	}

	if err := sanitizeMateriBankSoalPatch(&payload); err != nil {
		logger.Error(ctx, "failed update bank soal", "layer", "core.service", "op", "bank_soal.update", "err", err)
		return err
	}

	if err := r.repo.UpdateBankSoal(ctx, idBankSoal, payload); err != nil {
		logger.Error(ctx, "failed update bank soal", "layer", "core.service", "op", "bank_soal.update", "err", err)
		return err
	}
	return nil
}
