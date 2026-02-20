package pengumuman_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	pengumuman_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type GetPengumumanService struct {
	repo pengumuman_repo.PengumumanRepo
}

func NewGetPengumumanService(repo pengumuman_repo.PengumumanRepo) *GetPengumumanService {
	return &GetPengumumanService{
		repo: repo,
	}
}

func(r *GetPengumumanService)GetPengumumanActiveService(ctx context.Context)([]pengumuman.Pengumuman, error) {
	logger := corelog.FromContext(ctx)

	pengumumanActive, err := r.repo.GetPengumumanActive(ctx)
	if err != nil {
		logger.Error(ctx, "failed getting pengumuman", "layer", "core.service", "op", "pengumuman.get_pengumuman.GetPengumumanActiveService", "err", err)
		return nil, err
	}
	return pengumumanActive, nil
}

func(r *GetPengumumanService)GetPengumumanNonActiveService(ctx context.Context)([]pengumuman.Pengumuman, error) {
	logger := corelog.FromContext(ctx)

	pengumumanNonActive, err := r.repo.GetPengumumanNonActive(ctx)
	if err != nil {
		logger.Error(ctx, "failed getting pengumuman", "layer", "core.service", "op", "pengumuman.get_pengumuman.GetPengumumanNonActiveService", "err", err)
		return nil, err
	}
	return pengumumanNonActive, nil
}

func(r *GetPengumumanService)GetPengumumanIncomingService(ctx context.Context)([]pengumuman.Pengumuman, error) {
	logger := corelog.FromContext(ctx)

	pengumumanIncoming, err := r.repo.GetPengumumanIncoming(ctx)
	if err != nil {
		logger.Error(ctx, "failed getting pengumuman", "layer", "core.service", "op", "pengumuman.get_pengumuman.GetPengumumanIncomingService", "err", err)
		return nil, err
	}
	return pengumumanIncoming, nil
}

func(r *GetPengumumanService)GetPengumumanByIdService(ctx context.Context, idPengumuman pengumuman.ID)(pengumuman.Pengumuman, error) {
	logger := corelog.FromContext(ctx)

	if idPengumuman <= 0 {
		logger.Error(ctx, "failed getting pengumuman", "layer", "core.service", "op", "pengumuman.get_pengumuman.GetPengumumanByIdService", "err", coreerror.ErrMissingId)
		return pengumuman.Pengumuman{}, coreerror.ErrMissingId
	}

	item, err := r.repo.GetPengumumanById(ctx,idPengumuman)
	if err != nil {
		logger.Error(ctx, "failed getting pengumuman", "layer", "core.service", "op", "pengumuman.get_pengumuman.GetPengumumanByIdService", "err", err)
		return pengumuman.Pengumuman{}, err
	}
	return item, nil
}
