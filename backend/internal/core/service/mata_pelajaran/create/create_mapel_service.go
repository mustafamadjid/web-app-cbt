package matapelajaran_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	matapelajaran_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/mata_pelajaran"
)

type CreateMapelRepo struct {
	mapelRepo matapelajaran_repo.MataPelajaranRepository
}

func NewMapelService(mapelRepo matapelajaran_repo.MataPelajaranRepository) *CreateMapelRepo {
	return &CreateMapelRepo{
		mapelRepo: mapelRepo,
	}
}

func (r *CreateMapelRepo) CreateMapelService(ctx context.Context, mapel matapelajaran.MataPelajaran) error {
	logger := corelog.FromContext(ctx)

	mapel = sanitizeMapel(mapel)

	if err := validateMapelCreateInput(mapel); err != nil {
		return err
	}

	exist, err := r.mapelRepo.ExistKodeMapel(ctx, mapel.KodeMapel)
	if err != nil {
		logger.Error(ctx, "failed creating mapel", "layer", "core.service", "op", "matapelajaran.create_mapel.existKodeMapel", "err", err)
		return err
	}

	if exist {
		return coreerror.ErrKodeMapelExist
	}

	if err := r.mapelRepo.CreateMapel(ctx, mapel); err != nil {
		logger.Error(ctx, "failed creating mapel", "layer", "core.service", "op", "matapelajaran.create_mapel.CreateMapel", "err", err)
		return err
	}

	return nil
}
