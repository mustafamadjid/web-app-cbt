package ruangujian_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ruangujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateRuangUjianService struct {
	ruangRepo ruangujian_repo.RuangUjianRepo
}

func NewUpdateRuangUjianService(ruangRepo ruangujian_repo.RuangUjianRepo) *UpdateRuangUjianService {
	return &UpdateRuangUjianService{ruangRepo: ruangRepo}
}

func (r *UpdateRuangUjianService) UpdateRuangUjian(ctx context.Context, idRuangan int, ruangUjian updatepatch.UpdateRuangUjianPatch) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdateRuangUjianID(idRuangan); err != nil {
		logger.Error(ctx, "failed update ruang ujian", "layer", "core.service", "op", "ruangujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateRuangUjianPatch(ruangUjian); err != nil {
		logger.Error(ctx, "failed update ruang ujian", "layer", "core.service", "op", "ruangujian.update", "err", err)
		return err
	}

	sanitizeKodeRuangPatch(&ruangUjian)
	if isEmptyKodeRuangPatch(ruangUjian) {
		logger.Error(ctx, "failed updating ruangujian", "layer", "core.service", "op", "ruangujian.update_ruang_ujian.KodeRuang", "err", coreerror.ErrMissingField)
	}

	sanitizeNamaRuangPatch(&ruangUjian)
	if isEmptyNamaRuangPatch(ruangUjian) {
		logger.Error(ctx, "failed updating ruangujian", "layer", "core.service", "op", "ruangujian.update_ruang_ujian.NamaRuangan", "err", coreerror.ErrMissingField)
	}

	if ruangUjian.KodeRuang != nil {
		existKode, err := r.ruangRepo.ExistByKodeRuang(ctx, *ruangUjian.KodeRuang)
		if err != nil {
			logger.Error(ctx, "failed update ruang ujian", "layer", "core.service", "op", "ruangujian.update", "err", err)
			return err
		}

		if existKode {
			return coreerror.ErrKodeRuangUjianExist
		}
	}

	if err := r.ruangRepo.UpdateRuangUjian(ctx, idRuangan, ruangUjian); err != nil {
		logger.Error(ctx, "failed update ruang ujian", "layer", "core.service", "op", "ruangujian.update", "err", err)
		return err
	}
	return nil
}

// -----------------------
// Sanitizer and validator
// -----------------------

func validateUpdateRuangUjianID(idRuangan int) error {
	if idRuangan <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func validateUpdateRuangUjianPatch(ruangUjian updatepatch.UpdateRuangUjianPatch) error {
	if ruangUjian.KodeRuang == nil && ruangUjian.NamaRuang == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}

func sanitizeKodeRuangPatch(ruangUjian *updatepatch.UpdateRuangUjianPatch) {
	if ruangUjian.KodeRuang == nil {
		return
	}

	kodeRuang := strings.TrimSpace(*ruangUjian.KodeRuang)
	kodeRuang = strings.ToUpper(kodeRuang)
	ruangUjian.KodeRuang = &kodeRuang
}

func sanitizeNamaRuangPatch(ruangUjian *updatepatch.UpdateRuangUjianPatch) {
	if ruangUjian.NamaRuang == nil {
		return
	}

	namaRuang := strings.TrimSpace(*ruangUjian.NamaRuang)
	ruangUjian.NamaRuang = &namaRuang
}

func isEmptyKodeRuangPatch(ruangUjian updatepatch.UpdateRuangUjianPatch) bool {
	return ruangUjian.KodeRuang != nil && *ruangUjian.KodeRuang == ""
}

func isEmptyNamaRuangPatch(ruangUjian updatepatch.UpdateRuangUjianPatch) bool {
	return ruangUjian.NamaRuang != nil && *ruangUjian.NamaRuang == ""
}
