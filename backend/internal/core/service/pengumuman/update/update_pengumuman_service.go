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

// -----------------------
// Sanitizer and validator
// -----------------------

func validateUpdatePengumumanID(idPengumuman pengumuman.ID) error {
	if idPengumuman <= 0 {
		return coreerror.ErrMissingId
	}

	return nil
}

func validatePengumumanUpdateUserID(payload updatepatch.PengumumanUpdatePatch) error {
	if payload.IdPengguna != nil && *payload.IdPengguna <= 0 {
		return coreerror.ErrMissingId
	}

	return nil
}

func validatePengumumanUpdatePatch(payload updatepatch.PengumumanUpdatePatch) error {
	if payload.JudulPengumuman == nil &&
		payload.IsiPengumuman == nil &&
		payload.TanggalRilisPengumuman == nil &&
		payload.TanggalSelesaiPengumuman == nil &&
		payload.DokumenPengumuman == nil {
		return coreerror.ErrNoFieldToUpdate
	}

	return nil
}

func sanitizeJudulPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.JudulPengumuman == nil {
		return nil
	}

	judulPengumuman := strings.TrimSpace(*payload.JudulPengumuman)
	if judulPengumuman == "" {
		return coreerror.ErrMissingField
	}
	payload.JudulPengumuman = &judulPengumuman
	return nil
}

func sanitizeIsiPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.IsiPengumuman == nil {
		return nil
	}

	isiPengumuman := strings.TrimSpace(*payload.IsiPengumuman)
	if isiPengumuman == "" {
		return coreerror.ErrMissingField
	}
	payload.IsiPengumuman = &isiPengumuman
	return nil
}

func sanitizeTanggalRilisPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.TanggalRilisPengumuman == nil {
		return nil
	}

	tanggalRilisPengumuman := strings.TrimSpace(*payload.TanggalRilisPengumuman)
	if tanggalRilisPengumuman == "" {
		return coreerror.ErrMissingField
	}

	if err := pengumuman_service.ValidateDate(tanggalRilisPengumuman); err != nil {
		return err
	}

	payload.TanggalRilisPengumuman = &tanggalRilisPengumuman
	return nil
}

func sanitizeTanggalSelesaiPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.TanggalSelesaiPengumuman == nil {
		return nil
	}

	tanggalSelesaiPengumuman := strings.TrimSpace(*payload.TanggalSelesaiPengumuman)
	if tanggalSelesaiPengumuman == "" {
		return coreerror.ErrMissingField
	}

	if err := pengumuman_service.ValidateDate(tanggalSelesaiPengumuman); err != nil {
		return err
	}

	payload.TanggalSelesaiPengumuman = &tanggalSelesaiPengumuman
	return nil
}

func sanitizeDokumenPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.DokumenPengumuman == nil {
		return nil
	}

	dokumenPengumuman := strings.TrimSpace(*payload.DokumenPengumuman)
	if dokumenPengumuman == "" {
		return coreerror.ErrMissingField
	}

	payload.DokumenPengumuman = &dokumenPengumuman
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
