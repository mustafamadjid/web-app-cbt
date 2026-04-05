package gradingujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	grading_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/grading"
)

type StatistikUjianService struct {
	gradingRepo grading_repo.GradingUjianRepository
}


func NewStatistikUjianService(gradingRepo grading_repo.GradingUjianRepository) *StatistikUjianService {
	return &StatistikUjianService{gradingRepo: gradingRepo}
}

func(r *StatistikUjianService) StatistikUjian(ctx context.Context, idAttempt int) error {
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	if err := r.gradingRepo.UpsertToStatistikUjian(ctx,ujian.ID(idAttempt));
	err != nil {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	return nil
} 