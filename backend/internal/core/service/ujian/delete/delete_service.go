package ujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type DeleteUjianService struct {
	repo ujian_repo.UjianRepository
}

func NewDeleteUjianService(repo ujian_repo.UjianRepository) *DeleteUjianService {
	return &DeleteUjianService{repo: repo}
}

func (r *DeleteUjianService) DeleteUjianService(ctx context.Context, id ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if id <= 0 {
		logger.Error(ctx, "failed deleting ujian",
			"layer", "core.service",
			"op", "ujian.delete_ujian.DeleteUjianService",
			"err", coreerror.ErrMissingId,
		)
		return coreerror.ErrMissingId
	}

	if err := r.repo.DeleteUjian(ctx, id); err != nil {
		logger.Error(ctx, "failed deleting ujian",
			"layer", "core.service",
			"op", "ujian.delete_ujian.DeleteUjianService",
			"err", err,
		)
		return err
	}

	return nil
}

func (r *DeleteUjianService) DeleteJadwalUjianService(ctx context.Context, id ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if id <= 0 {
		logger.Error(ctx, "failed deleting jadwal ujian",
			"layer", "core.service",
			"op", "ujian.delete_jadwal.DeleteJadwalUjianService",
			"err", coreerror.ErrMissingId,
		)
		return coreerror.ErrMissingId
	}

	if err := r.repo.DeleteJadwalUjian(ctx, id); err != nil {
		logger.Error(ctx, "failed deleting jadwal ujian",
			"layer", "core.service",
			"op", "ujian.delete_jadwal.DeleteJadwalUjianService",
			"err", err,
		)
		return err
	}

	return nil
}

// Implementasi untuk yang dikomen:

func (r *DeleteUjianService) DeletePesertaUjianService(ctx context.Context, id ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if id <= 0 {
		logger.Error(ctx, "failed deleting peserta ujian",
			"layer", "core.service",
			"op", "ujian.delete_peserta.DeletePesertaUjianService",
			"err", coreerror.ErrMissingId,
		)
		return coreerror.ErrMissingId
	}

	if err := r.repo.DeletePesertaUjian(ctx, id); err != nil {
		logger.Error(ctx, "failed deleting peserta ujian",
			"layer", "core.service",
			"op", "ujian.delete_peserta.DeletePesertaUjianService",
			"err", err,
		)
		return err
	}

	return nil
}

func (r *DeleteUjianService) DeleteJawabanUjianSiswaService(ctx context.Context, id ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if id <= 0 {
		logger.Error(ctx, "failed deleting jawaban ujian siswa",
			"layer", "core.service",
			"op", "ujian.delete_jawaban.DeleteJawabanUjianSiswaService",
			"err", coreerror.ErrMissingId,
		)
		return coreerror.ErrMissingId
	}

	if err := r.repo.DeleteJawabanUjianSiswa(ctx, id); err != nil {
		logger.Error(ctx, "failed deleting jawaban ujian siswa",
			"layer", "core.service",
			"op", "ujian.delete_jawaban.DeleteJawabanUjianSiswaService",
			"err", err,
		)
		return err
	}

	return nil
}
