package pengumuman_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	pengumuman_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/pengumuman"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/date_validation"
)

type CreatePengumumanService struct {
	repo pengumuman_repo.PengumumanRepo
}

func NewCreatePengumumanRepo(repo pengumuman_repo.PengumumanRepo) *CreatePengumumanService {
	return &CreatePengumumanService{
		repo: repo,
	}
}

func(r *CreatePengumumanService)CreatePengumuman(ctx context.Context, pengumuman pengumuman.Pengumuman) error {
	logger := corelog.FromContext(ctx)

	pengumuman.IsiPengumuman = strings.TrimSpace(pengumuman.IsiPengumuman)
	pengumuman.JudulPengumuman = strings.TrimSpace(pengumuman.JudulPengumuman)
	pengumuman.TanggalRilisPengumuman = strings.TrimSpace(pengumuman.TanggalRilisPengumuman)
	pengumuman.TanggalSelesaiPengumuman = strings.TrimSpace(pengumuman.TanggalSelesaiPengumuman)
	pengumuman.DokumenPengumuman = strings.TrimSpace(pengumuman.DokumenPengumuman)

	if pengumuman.IdPengguna <= 0 {
		logger.Error(ctx, "failed creating pengumuman", "layer", "core.service", "op", "pengumuman.create_pengumuman.CreatePengumuman", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}
	
	if err := pengumuman_service.ValidateDate(pengumuman.TanggalRilisPengumuman); err != nil {
		logger.Error(ctx, "failed creating pengumuman", "layer", "core.service", "op", "pengumuman.create_pengumuman.ValidateDate", "err", err)
		return err
	}

	if err := pengumuman_service.ValidateDate(pengumuman.TanggalSelesaiPengumuman);err != nil {
		logger.Error(ctx, "failed creating pengumuman", "layer", "core.service", "op", "pengumuman.create_pengumuman.ValidateDate", "err", err)
		return err
	}

	if err := r.repo.CreatePengumuman(ctx,pengumuman);err != nil {
		logger.Error(ctx, "failed creating pengumuman", "layer", "core.service", "op", "pengumuman.create_pengumuman.CreatePengumuman", "err", err)
		return err
	}

	return nil

}

