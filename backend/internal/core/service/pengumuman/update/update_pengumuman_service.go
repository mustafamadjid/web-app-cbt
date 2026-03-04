package pengumuman_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	delete_file_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/delete_file_system"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	pengumuman_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdatePengumumanService struct {
	repo       pengumuman_repo.PengumumanRepo
	deleteFile delete_file_repo.DeleteFileRepo
}

func NewUpdatePengumumanService(repo pengumuman_repo.PengumumanRepo, deleteFile delete_file_repo.DeleteFileRepo) *UpdatePengumumanService {
	return &UpdatePengumumanService{
		repo:       repo,
		deleteFile: deleteFile,
	}
}

func (r *UpdatePengumumanService) UpdatePengumumanService(ctx context.Context, idPengumuman pengumuman.ID, pengumuman updatepatch.PengumumanUpdatePatch) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdatePengumumanID(idPengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validatePengumumanUpdateUserID(pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validatePengumumanUpdatePatch(pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := sanitizeJudulPengumumanPatch(&pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := sanitizeIsiPengumumanPatch(&pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := sanitizeTanggalRilisPengumumanPatch(&pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := sanitizeTanggalSelesaiPengumumanPatch(&pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := sanitizeDokumenPengumumanPatch(&pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := r.deleteOldDokumenPengumumanIfNeeded(ctx, idPengumuman, pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := r.repo.UpdatePengumuman(ctx, idPengumuman, pengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}
	return nil
}

func (r *UpdatePengumumanService) deleteOldDokumenPengumumanIfNeeded(ctx context.Context, idPengumuman pengumuman.ID, payload updatepatch.PengumumanUpdatePatch) error {
	if payload.DokumenPengumuman == nil {
		return nil
	}

	pengumumanData, err := r.repo.GetPengumumanById(ctx, idPengumuman)
	if err != nil {
		return err
	}

	if err := r.deleteFile.DeleteFile(ctx, pengumumanData.DokumenPengumuman); err != nil {
		return err
	}

	return nil
}
