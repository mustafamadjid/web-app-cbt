package bank_soal_service

import (
	"context"
	"strings"

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

// -----------------------
// Sanitizer and validator
// -----------------------

func validateUpdateBankSoalID(idBankSoal bank_soal.ID) error {
	if idBankSoal <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func validateUpdateBankSoalPatch(payload updatepatch.UpdateBankSoalPatch) error {
	if payload.IdMapel == nil &&
		payload.IdKelas == nil &&
		payload.IdPengguna == nil &&
		payload.NamaBankSoal == nil &&
		payload.Deskripsi == nil &&
		payload.Materi == nil {
		return coreerror.ErrNoFieldToUpdate
	}

	return nil
}

func validateUpdateBankSoalPatchID(payload updatepatch.UpdateBankSoalPatch) error {
	if payload.IdKelas != nil && *payload.IdKelas <= 0 {
		return coreerror.ErrMissingId
	}

	if payload.IdPengguna != nil && *payload.IdPengguna <= 0 {
		return coreerror.ErrMissingId
	}

	if payload.IdMapel != nil && *payload.IdMapel <= 0 {
		return coreerror.ErrMissingId
	}

	return nil
}

func sanitizeNamaBankSoalPatch(payload *updatepatch.UpdateBankSoalPatch) error {
	if payload.NamaBankSoal == nil {
		return nil
	}

	namaBankSoal := strings.TrimSpace(*payload.NamaBankSoal)
	if namaBankSoal == "" {
		return coreerror.ErrMissingField
	}

	payload.NamaBankSoal = &namaBankSoal
	return nil
}

func sanitizeDeskripsiBankSoalPatch(payload *updatepatch.UpdateBankSoalPatch) error {
	if payload.Deskripsi == nil {
		return nil
	}

	deskripsi := strings.TrimSpace(*payload.Deskripsi)
	if deskripsi == "" {
		return coreerror.ErrMissingField
	}

	payload.Deskripsi = &deskripsi
	return nil
}

func sanitizeMateriBankSoalPatch(payload *updatepatch.UpdateBankSoalPatch) error {
	if payload.Materi == nil {
		return nil
	}

	materi := strings.TrimSpace(*payload.Materi)
	if materi == "" {
		return coreerror.ErrMissingField
	}

	payload.Materi = &materi
	return nil
}
