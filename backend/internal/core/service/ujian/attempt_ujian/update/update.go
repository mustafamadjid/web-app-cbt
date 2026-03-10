package attemptujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type updateAttemptUjianRepository interface {
	UpdateAttemptUjian(ctx context.Context, idAttempt ujian.ID, data updatepatch.UpdateAttemptUjianPatch) error
}

type UpdateAttemptUjianService struct {
	repo updateAttemptUjianRepository
}

func NewUpdateAttemptUjianService(repo updateAttemptUjianRepository) *UpdateAttemptUjianService {
	return &UpdateAttemptUjianService{repo: repo}
}

func (r *UpdateAttemptUjianService) UpdateAttemptUjian(ctx context.Context, idAttempt ujian.ID, payload updatepatch.UpdateAttemptUjianPatch) error {
	logger := corelog.FromContext(ctx)

	if err := sanitizeUpdateAttemptStatusPatch(&payload); err != nil {
		logger.Error(ctx, "failed update attempt ujian", "layer", "core.service", "op", "ujian.attempt.update", "err", err)
		return err
	}

	if err := validateUpdateAttemptUjian(idAttempt, payload); err != nil {
		logger.Error(ctx, "failed update attempt ujian", "layer", "core.service", "op", "ujian.attempt.update", "err", err)
		return err
	}

	if err := r.repo.UpdateAttemptUjian(ctx, idAttempt, payload); err != nil {
		logger.Error(ctx, "failed update attempt ujian", "layer", "core.service", "op", "ujian.attempt.update", "err", err)
		return err
	}

	return nil
}
