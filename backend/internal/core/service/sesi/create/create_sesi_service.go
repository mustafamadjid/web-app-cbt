package sesi_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
)

type CreateSesiService struct {
	sesiRepo sesi_repo.SesiRepository
}

func NewCreateSesiService(sesiRepo sesi_repo.SesiRepository) *CreateSesiService {
	return &CreateSesiService{sesiRepo: sesiRepo}
}

func (r *CreateSesiService) CreateSesiService(ctx context.Context, sesi sesi.Sesi) error {
	logger := corelog.FromContext(ctx)

	sesi = sanitizeCreateSesi(sesi)

	if err := validateCreateNamaSesi(sesi); err != nil {
		logger.Error(ctx, "failed create sesi", "layer", "core.service", "op", "sesi.create", "err", coreerror.ErrMissingField)
		return err
	}

	if err := validateCreateKodeSesi(sesi); err != nil {
		logger.Error(ctx, "failed create sesi", "layer", "core.service", "op", "sesi.create", "err", coreerror.ErrMissingField)
		return err
	}

	exist, err := r.sesiRepo.ExistByKodeSesi(ctx, sesi.KodeSesi)
	if err != nil {
		logger.Error(ctx, "failed check exist sesi", "layer", "core.service", "op", "sesi.create", "err", err)
		return err
	}

	if exist {
		logger.Error(ctx, "failed create sesi", "layer", "core.service", "op", "sesi.create", "err", err)
		return coreerror.ErrSesiUjianExist
	}

	if err := r.sesiRepo.CreateSesi(ctx, sesi); err != nil {
		logger.Error(ctx, "failed create sesi", "layer", "core.service", "op", "sesi.create", "err", err)
		return err
	}
	return nil
}

// -----------------------
// Sanitizer and validator
// -----------------------

func sanitizeCreateSesi(data sesi.Sesi) sesi.Sesi {
	data.NamaSesi = strings.TrimSpace(data.NamaSesi)
	data.KodeSesi = strings.TrimSpace(data.KodeSesi)
	data.KodeSesi = strings.ToUpper(data.KodeSesi)
	return data
}

func validateCreateNamaSesi(data sesi.Sesi) error {
	if data.NamaSesi == "" {
		return coreerror.ErrMissingField
	}
	return nil
}

func validateCreateKodeSesi(data sesi.Sesi) error {
	if data.KodeSesi == "" {
		return coreerror.ErrMissingField
	}
	return nil
}
