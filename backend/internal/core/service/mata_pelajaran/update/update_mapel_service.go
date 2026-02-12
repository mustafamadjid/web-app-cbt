package matapelajaran_service

import (
	"context"
	"errors"
	"strings"

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

func (r *UpdateMapelRepo)UpdateMapelService(ctx context.Context, idMapel int, mapel updatepatch.UpdateMapelPatch)error{
	logger := corelog.FromContext(ctx)

	if idMapel <= 0 {
		return coreerror.ErrMissingId
	}

	if mapel.IdKelas != nil {
		if *mapel.IdKelas <= 0 {
			logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.IdKelas", "err", coreerror.ErrMissingId)
			return coreerror.ErrMissingId
		}
	}

	if mapel.KodeMapel != nil {
		if *mapel.KodeMapel == "" {
			logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.KodeMapel", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}

		s := strings.TrimSpace(*mapel.KodeMapel)
		s = strings.ToUpper(*mapel.KodeMapel)
		mapel.KodeMapel = &s

		exist, err := r.mapelRepo.ExistKodeMapel(ctx,*mapel.KodeMapel)
		if err != nil {
			logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.existKodeMapel", "err", err)
			return err
		}

		if exist {
			return coreerror.ErrKodeMapelExist
		}
	}

	if mapel.NamaMapel != nil {
		if *mapel.NamaMapel == "" {
			logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.NamaMapel", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}

		s := strings.TrimSpace(*mapel.NamaMapel)
		mapel.NamaMapel = &s
	}

	if mapel.Deskripsi != nil {
		if *mapel.Deskripsi == "" {
			logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.Deskripsi", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}

		s := strings.TrimSpace(*mapel.Deskripsi)
		mapel.Deskripsi = &s
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