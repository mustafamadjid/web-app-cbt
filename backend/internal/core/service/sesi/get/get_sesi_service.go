package sesi_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	authsession "github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type GetSesiService struct {
	sesiRepo    sesi_repo.SesiRepository
	sessionRepo out.SessionRepository
}

func NewGetSesiService(sesiRepo sesi_repo.SesiRepository, sessionRepo out.SessionRepository) *GetSesiService {
	return &GetSesiService{sesiRepo: sesiRepo, sessionRepo: sessionRepo}
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

func (r *GetSesiService) GetAllActiveSessionService(ctx context.Context) ([]authsession.SessionWithUser, error) {
	logger := corelog.FromContext(ctx)

	items, err := r.sessionRepo.GetAllActiveSession(ctx)
	if err != nil {
		logger.Error(ctx, "failed get all active login session", "layer", "core.service", "op", "sesi.get_all_active_login_session", "err", err)
		return nil, err
	}

	return items, nil
}
