package matapelajaran_service

import (
	"context"
	"errors"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	matapelajaran_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateMapelRepo struct {
	mapelRepo matapelajaran_repo.MataPelajaranRepository
}

func NewUpdateMapelService(mapelRepo matapelajaran_repo.MataPelajaranRepository) *UpdateMapelRepo {
	return &UpdateMapelRepo{
		mapelRepo: mapelRepo,
	}
}

func (r *UpdateMapelRepo) UpdateMapelService(ctx context.Context, idMapel int, mapel updatepatch.UpdateMapelPatch) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdateMapelID(idMapel); err != nil {
		return err
	}

	if err := validateMapelPatch(mapel); err != nil {
		return err
	}

	if err := validateMapelIdKelasPatch(mapel); err != nil {
		logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.IdKelas", "err", err)
		return err
	}

	if err := sanitizeKodeMapelPatch(&mapel); err != nil {
		logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.KodeMapel", "err", err)
		return err
	}

	if err := r.validateKodeMapelUniqueness(ctx, mapel); err != nil {
		if errors.Is(err, coreerror.ErrKodeMapelExist) {
			return err
		}
		logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.existKodeMapel", "err", err)
		return err
	}

	if err := sanitizeNamaMapelPatch(&mapel); err != nil {
		logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.NamaMapel", "err", err)
		return err
	}

	if err := sanitizeDeskripsiMapelPatch(&mapel); err != nil {
		logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.Deskripsi", "err", err)
		return err
	}

	if err := r.mapelRepo.UpdateMapel(ctx, idMapel, mapel); err != nil {
		logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.UpdateMapel", "err", err)

		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			return coreerror.ErrNotFound
		default:
			return err
		}
	}

	return nil
}
