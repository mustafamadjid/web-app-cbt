package matapelajajaran_service

import (
	"context"
	"errors"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	matapelajaran_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/mata_pelajaran"
)

type DeleteMapelRepo struct {
	mapelRepo matapelajaran_repo.MataPelajaranRepository
}

func NewDeleteMapelService (mapelRepo matapelajaran_repo.MataPelajaranRepository) *DeleteMapelRepo {
	return &DeleteMapelRepo{
		mapelRepo: mapelRepo,
	}
}

func (r *DeleteMapelRepo) DeleteMapelService(ctx context.Context, idMapel int) error {
	logger := corelog.FromContext(ctx)

	if idMapel <= 0 {
		return coreerror.ErrMissingId
	}

	if err := r.mapelRepo.DeleteMapel(ctx, idMapel); err != nil {
		logger.Error(ctx, "failed deleting mapel", "layer", "core.service", "op", "matapelajaran.delete_mapel.DeleteMapel", "err", err)
		
		switch {
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			return coreerror.ErrDeleteRestricted
		default:
			return err
		}
	}
	return nil
}