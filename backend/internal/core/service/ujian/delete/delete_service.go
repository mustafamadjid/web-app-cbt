package ujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type DeleteUjianService struct {
	ujianRepo UjianRepository
}

func NewDeleteUjianService(ujianRepo UjianRepository) *DeleteUjianService {
	return &DeleteUjianService{
		ujianRepo: ujianRepo,
	}
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

	if err := r.ujianRepo.DeleteUjian(ctx, id); err != nil {
		logger.Error(ctx, "failed deleting ujian",
			"layer", "core.service",
			"op", "ujian.delete_ujian.DeleteUjianService",
			"err", err,
		)
		return err
	}

	return nil
}

type DeletePesertaUjianService struct {
	pesertaRepo PesertaUjianRepository
}

func NewDeletePesertaUjianService(pesertaRepo PesertaUjianRepository) *DeletePesertaUjianService {
	return &DeletePesertaUjianService{
		pesertaRepo: pesertaRepo,
	}
}

func (r *DeletePesertaUjianService) DeletePesertaUjianService(ctx context.Context, id ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if id <= 0 {
		logger.Error(ctx, "failed deleting peserta ujian",
			"layer", "core.service",
			"op", "ujian.delete_peserta.DeletePesertaUjianService",
			"err", coreerror.ErrMissingId,
		)
		return coreerror.ErrMissingId
	}

	if err := r.pesertaRepo.DeletePesertaUjian(ctx, id); err != nil {
		logger.Error(ctx, "failed deleting peserta ujian",
			"layer", "core.service",
			"op", "ujian.delete_peserta.DeletePesertaUjianService",
			"err", err,
		)
		return err
	}

	return nil
}

type DeleteJawabanUjianSiswaService struct {
	jawabanRepo JawabanUjianRepository
}

func NewDeleteJawabanUjianSiswaService(jawabanRepo JawabanUjianRepository) *DeleteJawabanUjianSiswaService {
	return &DeleteJawabanUjianSiswaService{
		jawabanRepo: jawabanRepo,
	}
}

func (r *DeleteJawabanUjianSiswaService) DeleteJawabanUjianSiswaService(ctx context.Context, id ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if id <= 0 {
		logger.Error(ctx, "failed deleting jawaban ujian siswa",
			"layer", "core.service",
			"op", "ujian.delete_jawaban.DeleteJawabanUjianSiswaService",
			"err", coreerror.ErrMissingId,
		)
		return coreerror.ErrMissingId
	}

	if err := r.jawabanRepo.DeleteJawabanUjianSiswa(ctx, id); err != nil {
		logger.Error(ctx, "failed deleting jawaban ujian siswa",
			"layer", "core.service",
			"op", "ujian.delete_jawaban.DeleteJawabanUjianSiswaService",
			"err", err,
		)
		return err
	}

	return nil
}
