package attemptujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type GetActiveAttemptUjianService struct {
	repo ujian_repo.UjianSiswaRepository
}

func NewGetActiveAttemptUjianService(repo ujian_repo.UjianSiswaRepository) *GetActiveAttemptUjianService {
	return &GetActiveAttemptUjianService{
		repo: repo,
	}
}

func (r *GetActiveAttemptUjianService) GetActiveAttemptUjian(ctx context.Context, idSiswa int, idJadwalUjian int) (ujian.AttemptUjian, error) {
	logger := corelog.FromContext(ctx)

	if idJadwalUjian <= 0 {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrMissingId)
		return ujian.AttemptUjian{}, coreerror.ErrMissingId
	}

	if idSiswa <= 0 {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrMissingId)
		return ujian.AttemptUjian{}, coreerror.ErrMissingId
	}

	item, err := r.repo.GetActiveUjianAttemptBySiswa(ctx, idSiswa, idJadwalUjian)
	if err != nil {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", err)
		return ujian.AttemptUjian{}, err
	}

	return item, nil

}
