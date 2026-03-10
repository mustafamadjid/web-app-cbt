package attemptujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type deleteAttemptUjianRepository interface {
	DeleteAttemptUjian(ctx context.Context, idAttempt ujian.ID) error
}

type DeleteAttemptUjianService struct {
	repo deleteAttemptUjianRepository
}

func NewDeleteAttemptUjianService(repo deleteAttemptUjianRepository) *DeleteAttemptUjianService {
	return &DeleteAttemptUjianService{repo: repo}
}

func (r *DeleteAttemptUjianService) DeleteAttemptUjian(ctx context.Context, idAttempt ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx, "failed delete attempt ujian", "layer", "core.service", "op", "ujian.attempt.delete", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	if err := r.repo.DeleteAttemptUjian(ctx, idAttempt); err != nil {
		logger.Error(ctx, "failed delete attempt ujian", "layer", "core.service", "op", "ujian.attempt.delete", "err", err)
		return err
	}

	return nil
}
