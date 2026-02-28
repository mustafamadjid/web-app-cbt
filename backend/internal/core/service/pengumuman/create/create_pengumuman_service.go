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

// -----------------------
// Sanitizer and validator
// -----------------------

func sanitizeCreatePengumuman(data pengumuman.Pengumuman) pengumuman.Pengumuman {
	data.IsiPengumuman = strings.TrimSpace(data.IsiPengumuman)
	data.JudulPengumuman = strings.TrimSpace(data.JudulPengumuman)
	data.TanggalRilisPengumuman = strings.TrimSpace(data.TanggalRilisPengumuman)
	data.TanggalSelesaiPengumuman = strings.TrimSpace(data.TanggalSelesaiPengumuman)
	data.DokumenPengumuman = strings.TrimSpace(data.DokumenPengumuman)
	return data
}

func validateCreatePengumumanID(data pengumuman.Pengumuman) error {
	if data.IdPengguna <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func validateTanggalRilisPengumuman(data pengumuman.Pengumuman) error {
	return pengumuman_service.ValidateDate(data.TanggalRilisPengumuman)
}

func validateTanggalSelesaiPengumuman(data pengumuman.Pengumuman) error {
	return pengumuman_service.ValidateDate(data.TanggalSelesaiPengumuman)
}
