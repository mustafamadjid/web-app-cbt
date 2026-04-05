package attemptujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type PesertaUjianSubmittedService struct {
	attemptRepo ujian_repo.AttemptUjianRepository
}

func NewPesertaUjianSubmittedService(attemptRepo ujian_repo.AttemptUjianRepository) *PesertaUjianSubmittedService {
	return &PesertaUjianSubmittedService{
		attemptRepo: attemptRepo,
	}
}

func(r *PesertaUjianSubmittedService) ListPesertaUjianSubmitted(ctx context.Context,idJadwalUjian int)([]ujian.PesertaUjianSubmitted,error){
	logger := corelog.FromContext(ctx)

	if idJadwalUjian <= 0 {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrMissingId)
		return []ujian.PesertaUjianSubmitted{}, coreerror.ErrMissingId
	}

	items, err := r.attemptRepo.ListPesertaUjianAttemptSubmittedByIdJadwalUjian(ctx,ujian.ID(idJadwalUjian))
	if err != nil {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", err)
		return []ujian.PesertaUjianSubmitted{}, err
	}

	return items, nil
}