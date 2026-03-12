package ujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateUjianService struct {
	ujianRepo ujian_repo.UjianRepository
}

func NewUpdateUjianService(ujianRepo ujian_repo.UjianRepository) *UpdateUjianService {
	return &UpdateUjianService{
		ujianRepo: ujianRepo,
	}
}

func (r *UpdateUjianService) UpdateUjianService(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdatePenjadwalanUjianPatch(payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateUjianID(id); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateUjianPatchID(payload.Ujian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validateUpdateJadwalUjianPatchID(payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", coreerror.ErrMissingId)
		return err
	}

	if err := sanitizeNamaUjianPatch(&payload.Ujian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := sanitizeDeskripsiUjianPatch(&payload.Ujian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := sanitizeStatusUjianPatch(&payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := sanitizeTokenJadwalUjianPatch(&payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateJadwalUjianStatus(payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateJadwalUjianToken(payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := validateUpdateJadwalUjianTime(payload.JadwalUjian); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	if err := r.ujianRepo.UpdateUjian(ctx, id, payload); err != nil {
		logger.Error(ctx, "failed update ujian", "layer", "core.service", "op", "ujian.update", "err", err)
		return err
	}

	return nil
}
