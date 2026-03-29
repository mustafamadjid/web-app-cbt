package siswaujian_service

import (
	"context"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type GetWaktuSelesaiService struct {
	repo ujian_repo.UjianSiswaRepository
}

func NewGetWaktuSelesaiService(repo ujian_repo.UjianSiswaRepository) *GetWaktuSelesaiService {
	return &GetWaktuSelesaiService{repo: repo}
}

func (r *GetWaktuSelesaiService) GetWaktuSelesai(ctx context.Context, idJadwalUjian int) (time.Time, error) {
	logger := corelog.FromContext(ctx)

	if idJadwalUjian <= 0 {
		logger.Error(ctx, "invalid id jadwal ujian", "layer", "service", "op", "get_waktu_selesai")
		return time.Time{}, coreerror.ErrMissingId
	}

	waktuSelesai, err := r.repo.GetWaktuSelesaiUjian(ctx, idJadwalUjian)
	if err != nil {
		logger.Error(ctx, "failed get waktu selesai", "layer", "service", "op", "get_waktu_selesai", "err", err)
		return time.Time{}, err
	}

	return waktuSelesai, nil
}
