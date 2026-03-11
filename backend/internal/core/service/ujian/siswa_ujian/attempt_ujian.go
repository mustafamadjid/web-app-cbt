package siswaujian_service

import (
	"context"
	"strings"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type AttemptUjianService struct {
	repo        ujian_repo.SiswaUjianChecker
	attemptRepo ujian_repo.AttemptUjianRepository
}

func NewAttemptUjianService(repo ujian_repo.SiswaUjianChecker, attemptRepo ujian_repo.AttemptUjianRepository) *AttemptUjianService {
	return &AttemptUjianService{
		repo:        repo,
		attemptRepo: attemptRepo,
	}
}

func (r *AttemptUjianService) AttemptUjian(ctx context.Context, idSiswa int, idJadwalUjian int, tokenUjian string, waktuAttempt time.Time) error {
	logger := corelog.FromContext(ctx)

	if idSiswa <= 0 {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	if idJadwalUjian <= 0 {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	tokenUjian = strings.ToUpper(strings.TrimSpace(tokenUjian))
	if tokenUjian == "" {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrMissingTokenUjian)
		return coreerror.ErrMissingTokenUjian
	}

	if waktuAttempt.IsZero() {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrMissingId)
		return coreerror.ErrTimeEmpty
	}

	pesertaValid, idPesertaUjian, err := r.repo.CheckValidSiswaInPesertaUjianById(ctx, idSiswa, idJadwalUjian)
	if err != nil {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", err)
		return err
	}

	if !pesertaValid {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrPesertaNotAllowedToAttemptJadwal)
		return coreerror.ErrPesertaNotAllowedToAttemptJadwal
	}

	deadlineUjian, err := r.repo.GetDeadlineUjian(ctx, idJadwalUjian)
	if err != nil {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", err)
		return err
	}

	if !waktuAttempt.Before(deadlineUjian) && !waktuAttempt.Equal(deadlineUjian) {
		return coreerror.ErrWaktuAttemptPesertaInvalid
	}

	tokenValid, err := r.repo.CheckTokenUjian(ctx, tokenUjian)
	if err != nil {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", err)
		return err
	}

	if !tokenValid {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", coreerror.ErrTokenUjianInvalid)
		return coreerror.ErrTokenUjianInvalid
	}

	dataAttempt := ujian.AttemptUjian{
		IdPesertaUjian: ujian.ID(idPesertaUjian),
		WaktuMulai:     &waktuAttempt,
		DeadlineAt:     &deadlineUjian,
	}

	if err := r.attemptRepo.CreateAttemptUjian(ctx, dataAttempt); err != nil {
		logger.Error(ctx, "failed attempt ujian", "layer", "core.service", "op", "ujian.attempt.create", "err", err)
		return err
	}

	return nil

}
