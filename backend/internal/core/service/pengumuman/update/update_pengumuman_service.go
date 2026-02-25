package pengumuman_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	delete_file_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/delete_file_system"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	pengumuman_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/date_validation"
)

type UpdatePengumumanService struct {
	repo pengumuman_repo.PengumumanRepo
	deleteFile delete_file_repo.DeleteFileRepo
}

func NewUpdatePengumumanService(repo pengumuman_repo.PengumumanRepo, deleteFile delete_file_repo.DeleteFileRepo) *UpdatePengumumanService {
	return &UpdatePengumumanService{
		repo: repo,
		deleteFile: deleteFile,
	}
}

func(r *UpdatePengumumanService)UpdatePengumumanService(ctx context.Context, idPengumuman pengumuman.ID, pengumuman updatepatch.PengumumanUpdatePatch) error {
	logger := corelog.FromContext(ctx)

	if idPengumuman <= 0 {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	if pengumuman.IdPengguna != nil {
		if *pengumuman.IdPengguna <= 0 {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingId)
			return coreerror.ErrMissingId
		}
	}

	if pengumuman.JudulPengumuman != nil {
		judulPengumuman := strings.TrimSpace(*pengumuman.JudulPengumuman)
		if judulPengumuman == "" {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}
		pengumuman.JudulPengumuman = &judulPengumuman
	}

	if pengumuman.IsiPengumuman != nil {
		isiPengumuman := strings.TrimSpace(*pengumuman.IsiPengumuman)
		if isiPengumuman == "" {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}
		pengumuman.IsiPengumuman = &isiPengumuman
	}

	if pengumuman.TanggalRilisPengumuman != nil {
		tanggalRilisPengumuman := strings.TrimSpace(*pengumuman.TanggalRilisPengumuman)
		if tanggalRilisPengumuman == "" {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}

		if err:= pengumuman_service.ValidateDate(tanggalRilisPengumuman); err != nil {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
			return err
		}

		pengumuman.TanggalRilisPengumuman = &tanggalRilisPengumuman
	}

	if pengumuman.TanggalSelesaiPengumuman != nil {
		tanggalSelesaiPengumuman := strings.TrimSpace(*pengumuman.TanggalSelesaiPengumuman)
		if tanggalSelesaiPengumuman == "" {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}

		if err:= pengumuman_service.ValidateDate(tanggalSelesaiPengumuman); err != nil {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
			return err
		}

		pengumuman.TanggalSelesaiPengumuman = &tanggalSelesaiPengumuman
	}

	if pengumuman.DokumenPengumuman != nil {
		dokumenPengumuman := strings.TrimSpace(*pengumuman.DokumenPengumuman)
		if dokumenPengumuman == "" {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}

		pengumumanData,err := r.repo.GetPengumumanById(ctx,idPengumuman)
		if err != nil {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
			return err
		}

		if err := r.deleteFile.DeleteFile(ctx,pengumumanData.DokumenPengumuman); err != nil {
			logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.delete_file_foto", "pengumuman_id", idPengumuman, "err", err)
			return err
		}

		pengumuman.DokumenPengumuman = &dokumenPengumuman
	}

	if err := r.repo.UpdatePengumuman(ctx,idPengumuman,pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}
	return nil
}