package statisktikujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	statistikujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/statistik_ujian"
)

type StatistikUjianService struct {
	statistikRepo statistikujian_repo.StatistikUjianRepository
}

func NewStatistikUjianService(statistikRepo statistikujian_repo.StatistikUjianRepository) *StatistikUjianService {
	return &StatistikUjianService{statistikRepo: statistikRepo}
}

func(r *StatistikUjianService)GetStatistikUjian(ctx context.Context, idJadwal int) (ujian.StatistikUjian, error) {
	logger := corelog.FromContext(ctx)

	if idJadwal <= 0 {
		logger.Error(ctx, "failed get statistik ujian", "layer", "core.service", "op", "ujian.statistik", "err", coreerror.ErrMissingId)
		return ujian.StatistikUjian{}, coreerror.ErrMissingId
	}

	items,err := r.statistikRepo.GetStatistikUjianByIdJadwal(ctx,ujian.ID(idJadwal))
	if err != nil {
		logger.Error(ctx, "failed get statistik ujian", "layer", "core.service", "op", "ujian.statistik", "err", err)
		return ujian.StatistikUjian{}, err
	}

	return items, nil
}