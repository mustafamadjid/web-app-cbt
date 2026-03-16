package siswaujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type GetJawabanUjianService struct {
	repo ujian_repo.JawabanUjianRepository
}

func NewGetJawabanUjianService(repo ujian_repo.JawabanUjianRepository) *GetJawabanUjianService {
	return &GetJawabanUjianService{
		repo: repo,
	}
}

func (s *GetJawabanUjianService) GetJawabanUjianByAttemptId(ctx context.Context, idAttempt ujian.ID) ([]ujian.JawabanUjian, error) {
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx, "failed get jawaban ujian by attempt id", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id", "err", coreerror.ErrMissingId)
		return []ujian.JawabanUjian{}, coreerror.ErrMissingId
	}

	items, err := s.repo.GetJawabanUjianByAttemptId(ctx, idAttempt)
	if err != nil {
		logger.Error(ctx, "failed get jawaban ujian by attempt id", "layer", "core.service", "op", "ujian.jawaban.get_by_attempt_id", "err", err)
		return []ujian.JawabanUjian{}, err
	}

	return items, nil
}
