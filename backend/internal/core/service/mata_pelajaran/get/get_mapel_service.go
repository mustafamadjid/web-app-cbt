package matapelajaran_service

import (
	"context"
	"errors"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	matapelajaran_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/mata_pelajaran"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type GetMapelRepo struct {
	mapelRepo matapelajaran_repo.MataPelajaranRepository
}

func NewGetMapelService(mapelRepo matapelajaran_repo.MataPelajaranRepository) *GetMapelRepo {
	return &GetMapelRepo{
		mapelRepo: mapelRepo,
	}
}

func (r *GetMapelRepo) GetMapelService(ctx context.Context, filter query.ListMapelFilter) ([]matapelajaran.MataPelajaran, error) {
	logger := corelog.FromContext(ctx)

	var err error
	filter, err = sanitizeAndValidateListMapelFilter(filter)
	if err != nil {
		switch {
		case isNamaMapelError(err):
			logger.Error(ctx, "failed get mapel", "layer", "core.service", "op", "matapelajaran.get.NamaMapel", "err", coreerror.ErrMissingField)
		default:
			logger.Error(ctx, "failed get mapel", "layer", "core.service", "op", "matapelajaran.get.TingkatKelas", "err", coreerror.ErrMissingField)
		}
		return nil, coreerror.ErrInvalidInput
	}

	items, err := r.mapelRepo.GetMapel(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed get mapel", "layer", "core.service", "op", "matapelajaran.get", "err", err)
		return nil, err
	}

	return items, nil
}

func (r *GetMapelRepo) GetMapelById(ctx context.Context, idMapel int) (matapelajaran.MataPelajaran, error) {
	logger := corelog.FromContext(ctx)

	if err := validateMapelID(idMapel); err != nil {
		logger.Error(ctx, "failed get mapel by id", "layer", "core.service", "op", "matapelajaran.get_by_id", "err", coreerror.ErrMissingId)
		return matapelajaran.MataPelajaran{}, err
	}

	item, err := r.mapelRepo.GetMapelById(ctx, idMapel)
	if err != nil {
		logger.Error(ctx, "failed get mapel by id", "layer", "core.service", "op", "matapelajaran.get_by_id", "err", err)

		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			return matapelajaran.MataPelajaran{}, coreerror.ErrNotFound
		default:
			return matapelajaran.MataPelajaran{}, err
		}
	}

	return item, nil
}

var (
	errInvalidNamaMapelFilter = errors.New("invalid nama mapel")
	errInvalidTingkatKelas    = errors.New("invalid tingkat kelas")
)

func isNamaMapelError(err error) bool {
	return errors.Is(err, errInvalidNamaMapelFilter)
}
