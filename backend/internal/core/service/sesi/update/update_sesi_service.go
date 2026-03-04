package sesi_service

import (
	"context"
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

	if err := validateUpdateSesiPatch(sesi); err != nil {
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
