package gradingujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	grading_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/grading"
)

type EssayGradingUjianService struct {
	gradingRepo grading_repo.GradingUjianRepository
}

func NewEssayGradingUjianService(gradingRepo grading_repo.GradingUjianRepository) *EssayGradingUjianService {
	return &EssayGradingUjianService{gradingRepo: gradingRepo}
}

func (r *EssayGradingUjianService) EssayGrading(ctx context.Context, jawabanSiswa []ujian.JawabanUjian, gradedBy ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if gradedBy <= 0 {
		logger.Error(ctx, "failed essay grading", "layer", "core.service", "op", "ujian.grading", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	for _, itemJawaban := range jawabanSiswa {
		if itemJawaban.IdJawaban <= 0 {
			logger.Error(ctx, "failed essay grading", "layer", "core.service", "op", "ujian.grading", "err", coreerror.ErrMissingId)
			return coreerror.ErrMissingId
		}

		if itemJawaban.EssayIsBenar == nil {
			logger.Error(ctx, "failed essay grading", "layer", "core.service", "op", "ujian.grading", "err", coreerror.ErrMissingId)
			return coreerror.ErrMissingId
		}
	}

	if err := r.gradingRepo.UpdateAndGradingEssayUjian(ctx, jawabanSiswa, gradedBy); err != nil {
		logger.Error(ctx, "failed essay grading", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	return nil
}
