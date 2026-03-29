package attemptujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type SiswaUpdateAttemptUjianService struct {
	checker ujian_repo.SiswaUjianChecker
	updater *UpdateAttemptUjianService
}

func NewSiswaUpdateAttemptUjianService(checker ujian_repo.SiswaUjianChecker, updater *UpdateAttemptUjianService) *SiswaUpdateAttemptUjianService {
	return &SiswaUpdateAttemptUjianService{
		checker: checker,
		updater: updater,
	}
}

func (s *SiswaUpdateAttemptUjianService) UpdateAttemptUjian(ctx context.Context, idSiswa int, idAttempt ujian.ID, payload updatepatch.UpdateAttemptUjianPatch) error {
	logger := corelog.FromContext(ctx)

	if err := validateSiswaUpdateAttemptUjian(idSiswa, idAttempt, payload); err != nil {
		logger.Error(ctx, "failed update attempt ujian by siswa", "layer", "core.service", "op", "ujian.attempt.update_siswa", "err", err)
		return err
	}

	owned, err := s.checker.CheckAttemptOwnershipBySiswa(ctx, idSiswa, idAttempt)
	if err != nil {
		logger.Error(ctx, "failed update attempt ujian by siswa", "layer", "core.service", "op", "ujian.attempt.update_siswa", "err", err)
		return err
	}
	if !owned {
		logger.Error(ctx, "failed update attempt ujian by siswa", "layer", "core.service", "op", "ujian.attempt.update_siswa", "err", coreerror.ErrNotFound)
		return coreerror.ErrNotFound
	}

	if err := s.updater.UpdateAttemptUjian(ctx, idAttempt, payload); err != nil {
		logger.Error(ctx, "failed update attempt ujian by siswa", "layer", "core.service", "op", "ujian.attempt.update_siswa", "err", err)
		return err
	}

	return nil
}
