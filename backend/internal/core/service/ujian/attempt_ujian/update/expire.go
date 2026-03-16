package attemptujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type ExpireAttemptUjianService struct {
	updater *UpdateAttemptUjianService
}

func NewExpireAttemptUjianService(updater *UpdateAttemptUjianService) *ExpireAttemptUjianService {
	return &ExpireAttemptUjianService{updater: updater}
}

func (s *ExpireAttemptUjianService) ExpireAttemptUjian(ctx context.Context, idAttempt ujian.ID) error {
	logger := corelog.FromContext(ctx)

	status := ujian.ATTEMPT_EXPIRED
	patch := updatepatch.UpdateAttemptUjianPatch{
		StatusAttempt: &status,
	}

	if err := s.updater.UpdateAttemptUjian(ctx, idAttempt, patch); err != nil {
		logger.Error(ctx, "failed expire attempt ujian", "layer", "core.service", "op", "ujian.attempt.expire_admin", "err", err)
		return err
	}

	return nil
}
