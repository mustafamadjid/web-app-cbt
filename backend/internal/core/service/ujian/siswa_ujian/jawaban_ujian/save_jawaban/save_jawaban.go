package siswaujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type JawabanUjianService struct {
	repo ujian_repo.JawabanUjianRepository
}

func NewJawabanUjianService(repo ujian_repo.JawabanUjianRepository) *JawabanUjianService {
	return &JawabanUjianService{
		repo: repo,
	}
}

func(r *JawabanUjianService) SaveJawabanUjian(ctx context.Context, idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error {
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx, "failed save jawaban ujian", "layer", "core.service", "op", "ujian.jawaban.save", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	for _,item := range jawaban{
		if item.IdJawaban <= 0 {
			logger.Error(ctx, "failed save jawaban ujian", "layer", "core.service", "op", "ujian.jawaban.save", "err", coreerror.ErrMissingId)
			return coreerror.ErrMissingId
		}

		if item.IdPilihan == nil && item.JawabanEssay == nil {
			logger.Error(ctx, "failed save jawaban ujian", "layer", "core.service", "op", "ujian.jawaban.save", "err", coreerror.ErrMissingJawabanEssayAndPilgan)
			return coreerror.ErrMissingJawabanEssayAndPilgan
		}
	}

	if err := r.repo.SaveJawabanUjian(ctx, idAttempt, jawaban); err != nil {
		logger.Error(ctx, "failed save jawaban ujian", "layer", "core.service", "op", "ujian.jawaban.save", "err", err)
		return err
	}

	return nil
}