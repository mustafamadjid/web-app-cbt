package sesi_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type GetSesiService struct {
	sesiRepo sesi_repo.SesiRepository
}

func NewGetSesiService(sesiRepo sesi_repo.SesiRepository) *GetSesiService {
	return &GetSesiService{sesiRepo: sesiRepo}
}

func (r *GetSesiService) GetSesiService(ctx context.Context, filter query.ListSesiFilter) ([]sesi.Sesi, error) {
	logger := corelog.FromContext(ctx)

	filter = sanitizeListSesiFilter(filter)

	items, err := r.sesiRepo.GetSesi(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed get sesi", "layer", "core.service", "op", "sesi.get", "err", err)
		return nil, err
	}
	return items, nil
}

func (r *GetSesiService) GetSesiByIdService(ctx context.Context, idSesi int) (sesi.Sesi, error) {
	logger := corelog.FromContext(ctx)

	if err := validateSesiID(idSesi); err != nil {
		return sesi.Sesi{}, err
	}

	item, err := r.sesiRepo.GetSesiById(ctx, idSesi)
	if err != nil {
		logger.Error(ctx, "failed get sesi", "layer", "core.service", "op", "sesi.get_by_id", "err", err)
		return sesi.Sesi{}, err
	}
	return item, nil
}

func (r *GetSesiService) GetSesiByKodeService(ctx context.Context, kodeSesi string) (sesi.Sesi, error) {
	logger := corelog.FromContext(ctx)

	kodeSesi = sanitizeKodeSesi(kodeSesi)

	if err := validateKodeSesi(kodeSesi); err != nil {
		logger.Error(ctx, "failed get sesi", "layer", "core.service", "op", "sesi.get_by_kode", "err", coreerror.ErrMissingField)
		return sesi.Sesi{}, err
	}

	item, err := r.sesiRepo.GetSesiByKode(ctx, kodeSesi)
	if err != nil {
		logger.Error(ctx, "failed get sesi", "layer", "core.service", "op", "sesi.get_by_kode", "err", err)
		return sesi.Sesi{}, err
	}
	return item, nil
}

// -----------------------
// Sanitizer and validator
// -----------------------

func sanitizeListSesiFilter(filter query.ListSesiFilter) query.ListSesiFilter {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	if filter.Limit > 50 {
		filter.Limit = 50
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return filter
}

func validateSesiID(idSesi int) error {
	if idSesi <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func sanitizeKodeSesi(kodeSesi string) string {
	kodeSesi = strings.TrimSpace(kodeSesi)
	kodeSesi = strings.ToUpper(kodeSesi)
	return kodeSesi
}

func validateKodeSesi(kodeSesi string) error {
	if len(kodeSesi) == 0 || kodeSesi == "" {
		return coreerror.ErrMissingField
	}
	return nil
}
