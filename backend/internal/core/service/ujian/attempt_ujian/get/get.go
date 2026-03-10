package attemptujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type getAttemptUjianRepository interface {
	GetAttemptUjianById(ctx context.Context, idAttempt ujian.ID) (ujian.AttemptUjian, error)
}

type GetAttemptUjianService struct {
	repo getAttemptUjianRepository
}

func NewGetAttemptUjianService(repo getAttemptUjianRepository) *GetAttemptUjianService {
	return &GetAttemptUjianService{repo: repo}
}

func (r *GetAttemptUjianService) GetAttemptUjianById(ctx context.Context, idAttempt ujian.ID) (ujian.AttemptUjian, error) {
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx, "failed get attempt ujian", "layer", "core.service", "op", "ujian.attempt.get_by_id", "err", coreerror.ErrMissingId)
		return ujian.AttemptUjian{}, coreerror.ErrMissingId
	}

	item, err := r.repo.GetAttemptUjianById(ctx, idAttempt)
	if err != nil {
		logger.Error(ctx, "failed get attempt ujian", "layer", "core.service", "op", "ujian.attempt.get_by_id", "err", err)
		return ujian.AttemptUjian{}, err
	}

	return item, nil
}
