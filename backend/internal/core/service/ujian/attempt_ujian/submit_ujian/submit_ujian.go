package attemptujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type SubmitUjianService struct {
	repo    ujian_repo.AttemptUjianRepository
	checker ujian_repo.SiswaAttemptOwnershipChecker
}

func NewSubmitUjianService(repo ujian_repo.AttemptUjianRepository, checker ujian_repo.SiswaAttemptOwnershipChecker) *SubmitUjianService {
	return &SubmitUjianService{
		repo:    repo,
		checker: checker,
	}
}

func (r *SubmitUjianService) SubmitUjian(ctx context.Context, idAttempt ujian.ID, idSiswa int) error {
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx, "failed submit ujian", "layer", "core.service", "op", "ujian.attempt.submit", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	owned, err := r.checker.CheckAttemptOwnershipBySiswa(ctx, idSiswa, idAttempt)
	if err != nil {
		logger.Error(ctx, "failed update attempt ujian by siswa", "layer", "core.service", "op", "ujian.attempt.update_siswa", "err", err)
		return err
	}

	if !owned {
		logger.Error(ctx, "failed update attempt ujian by siswa", "layer", "core.service", "op", "ujian.attempt.update_siswa", "err", coreerror.ErrNotFound)
		return coreerror.ErrNotFound
	}

	if err := r.repo.SubmitAttemptUjian(ctx, idAttempt); err != nil {
		logger.Error(ctx, "failed submit ujian", "layer", "core.service", "op", "ujian.attempt.submit", "err", err)
		return err
	}

	return nil
}
