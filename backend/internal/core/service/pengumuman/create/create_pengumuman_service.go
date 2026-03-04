package pengumuman_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	pengumuman_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/pengumuman"
)

type CreatePengumumanService struct {
	repo pengumuman_repo.PengumumanRepo
}

func NewCreatePengumumanRepo(repo pengumuman_repo.PengumumanRepo) *CreatePengumumanService {
	return &CreatePengumumanService{
		repo: repo,
	}
}

func (r *CreatePengumumanService) CreatePengumuman(ctx context.Context, pengumuman pengumuman.Pengumuman) error {
	logger := corelog.FromContext(ctx)

	pengumuman = sanitizeCreatePengumuman(pengumuman)

	if err := validateCreatePengumumanID(pengumuman); err != nil {
		logger.Error(ctx, "failed creating pengumuman", "layer", "core.service", "op", "pengumuman.create_pengumuman.CreatePengumuman", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateTanggalRilisPengumuman(pengumuman); err != nil {
		logger.Error(ctx, "failed creating pengumuman", "layer", "core.service", "op", "pengumuman.create_pengumuman.ValidateDate", "err", err)
		return err
	}

	if err := validateTanggalSelesaiPengumuman(pengumuman); err != nil {
		logger.Error(ctx, "failed creating pengumuman", "layer", "core.service", "op", "pengumuman.create_pengumuman.ValidateDate", "err", err)
		return err
	}

	if err := r.repo.CreatePengumuman(ctx, pengumuman); err != nil {
		logger.Error(ctx, "failed creating pengumuman", "layer", "core.service", "op", "pengumuman.create_pengumuman.CreatePengumuman", "err", err)
		return err
	}

	return nil

}
