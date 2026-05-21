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

	mapel = sanitizeUpdateMapelPatch(mapel)

	if err := r.validateUpdateMapelPatch(ctx, idMapel, mapel); err != nil {
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
