package ruangujian_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ruangujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ruang_ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
)

type GetRuangUjianRepo struct {
	ruangRepo ruangujian_repo.RuangUjianRepo
}

func NewGetRuangUjianService(ruangRepo ruangujian_repo.RuangUjianRepo) *GetRuangUjianRepo {
	return &GetRuangUjianRepo{ruangRepo: ruangRepo}
}

func(r *GetRuangUjianRepo)GetRuangUjian(ctx context.Context, filter query.ListRuangUjianFilter)([]ruangujian.RuangUjian, error){
	logger := corelog.FromContext(ctx)

	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.Limit > 50 {
		filter.Limit = 50
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	filter.Search = strings.TrimSpace(filter.Search)

	items,err := r.ruangRepo.GetRuangUjian(ctx,filter)
	if err != nil {
		logger.Error(ctx,"failed get ruang ujian","layer","core.service","op","ruangujian.get","err",err)
		return nil,err
	}
	return items, nil
	
}

func(r *GetRuangUjianRepo)GetRuangUjianById(ctx context.Context, idRuangan int)(ruangujian.RuangUjian, error){
	logger := corelog.FromContext(ctx)

	if idRuangan <= 0 {
		logger.Error(ctx,"failed get ruang ujian","layer","core.service","op","ruangujian.get_by_id","err",coreerror.ErrMissingId)
		return ruangujian.RuangUjian{},coreerror.ErrMissingId
	}

	item,err := r.ruangRepo.GetRuangUjianById(ctx,idRuangan)
	if err != nil {
		logger.Error(ctx,"failed get ruang ujian","layer","core.service","op","ruangujian.get_by_id","err",err)
		return ruangujian.RuangUjian{},err
	}
	return item, nil
}

func(r *GetRuangUjianRepo)GetRuangUjianByKode(ctx context.Context, kodeRuang string)(ruangujian.RuangUjian, error){
	logger := corelog.FromContext(ctx)

	kodeRuang = strings.TrimSpace(kodeRuang)
	kodeRuang =strings.ToUpper(kodeRuang)

	if len(kodeRuang) == 0 || kodeRuang == "" {
		logger.Error(ctx,"failed get ruang ujian","layer","core.service","op","ruangujian.get_by_kode","err",coreerror.ErrMissingField)
		return ruangujian.RuangUjian{},coreerror.ErrMissingField
	}

	item,err := r.ruangRepo.GetRuangUjianByKode(ctx,kodeRuang)
	if err != nil {
		logger.Error(ctx,"failed get ruang ujian","layer","core.service","op","ruangujian.get_by_kode","err",err)
		return ruangujian.RuangUjian{},err
	}
	return item, nil
}