package attemptujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type CreateAttemptUjianService struct {
	repo ujian_repo.AttemptUjianRepository
}

func NewCreateAttemptUjianService(repo ujian_repo.AttemptUjianRepository) *CreateAttemptUjianService {
	return &CreateAttemptUjianService{
		repo: repo,
	}
}

func (r *CreateAttemptUjianService) CreateAttemptUjian(ctx context.Context, data ujian.AttemptUjian) error {
	logger := corelog.FromContext(ctx)

	data = sanitizeCreateAttemptUjian(data)

	if err := validateCreateAttemptUjian(data); err != nil {
		logger.Error(ctx, "failed create attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", err)
		return err
	}

	if err := r.repo.CreateAttemptUjian(ctx, data); err != nil {
		logger.Error(ctx, "failed create attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", err)
		return err
	}

	return nil
}
