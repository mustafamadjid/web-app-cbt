package sesi_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateSesiService struct {
	sesiRepo sesi_repo.SesiRepository
}

func NewUpdateSesiService(sesiRepo sesi_repo.SesiRepository) *UpdateSesiService {
	return &UpdateSesiService{sesiRepo: sesiRepo}
}

func (r *UpdateSesiService) UpdateSesiService(ctx context.Context, idSesi int, sesi updatepatch.UpdateSesiPatch) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdateSesiID(idSesi); err != nil {
		return err
	}

	if err := sanitizeNamaSesiPatch(&sesi); err != nil {
		return err
	}

	if err := sanitizeKodeSesiPatch(&sesi); err != nil {
		return err
	}

	if err := r.sesiRepo.UpdateSesi(ctx, idSesi, sesi); err != nil {
		logger.Error(ctx, "failed update sesi", "layer", "core.service", "op", "sesi.update", "err", err)
		return err
	}
	return nil
}

// -----------------------
// Sanitizer and validator
// -----------------------

func validateUpdateSesiID(idSesi int) error {
	if idSesi == 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func sanitizeNamaSesiPatch(sesi *updatepatch.UpdateSesiPatch) error {
	if sesi.NamaSesi == nil {
		return nil
	}

	namaSesi := strings.TrimSpace(*sesi.NamaSesi)
	if namaSesi == "" {
		return coreerror.ErrMissingField
	}

	sesi.NamaSesi = &namaSesi
	return nil
}

func sanitizeKodeSesiPatch(sesi *updatepatch.UpdateSesiPatch) error {
	if sesi.KodeSesi == nil {
		return nil
	}

	kodeSesi := strings.TrimSpace(*sesi.KodeSesi)
	if kodeSesi == "" {
		return coreerror.ErrMissingField
	}

	kodeSesi = strings.ToUpper(kodeSesi)
	sesi.KodeSesi = &kodeSesi
	return nil
}
