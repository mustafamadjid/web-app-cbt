package matapelajaran_service

import (
	"context"
	"errors"
	"strings"

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

func(r *GetMapelRepo)GetMapelService(ctx context.Context,filter query.ListMapelFilter)([]matapelajaran.MataPelajaran,error){
	logger := corelog.FromContext(ctx)

	filter.Search = strings.TrimSpace(filter.Search)
	
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.Limit > 50 {
		filter.Limit = 50
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	if filter.NamaMapel != nil {
		if *filter.NamaMapel == "" {
			logger.Error(ctx,"failed get mapel","layer","core.service","op","matapelajaran.get.NamaMapel","err",coreerror.ErrMissingField)
			return nil, coreerror.ErrInvalidInput
		}

		s := strings.TrimSpace(*filter.NamaMapel)
		filter.NamaMapel = &s
	}

	if filter.TingkatKelas != nil {
		if *filter.TingkatKelas <= 0 {
			logger.Error(ctx,"failed get mapel","layer","core.service","op","matapelajaran.get.TingkatKelas","err",coreerror.ErrMissingField)
			return nil, coreerror.ErrInvalidInput
		}
	}

	items, err := r.mapelRepo.GetMapel(ctx,filter)
	if err != nil {
		logger.Error(ctx,"failed get mapel","layer","core.service","op","matapelajaran.get","err",err)
		return nil,err
	}

	return items, nil
}

func (r *GetMapelRepo)GetMapelById(ctx context.Context, idMapel int) (matapelajaran.MataPelajaran, error){
	logger := corelog.FromContext(ctx)

	if idMapel <= 0 {
		return matapelajaran.MataPelajaran{}, coreerror.ErrMissingId
	}

	item, err := r.mapelRepo.GetMapelById(ctx,idMapel)
	if err != nil {
		logger.Error(ctx,"failed get mapel by id","layer","core.service","op","matapelajaran.get_by_id","err",err)
		
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			return matapelajaran.MataPelajaran{}, coreerror.ErrNotFound
		default:
			return matapelajaran.MataPelajaran{}, err
		}
	}

	return item, nil
}