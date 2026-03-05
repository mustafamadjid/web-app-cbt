package ujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type CreateUjianService struct {
	ujianRepo UjianRepository
}

func NewCreateUjianService(ujianRepo UjianRepository) *CreateUjianService {
	return &CreateUjianService{
		ujianRepo: ujianRepo,
	}
}

func (r *CreateUjianService) CreateUjianService(ctx context.Context, data ujian.PenjadwalanUjian) error {
	logger := corelog.FromContext(ctx)

	data = sanitizeCreateUjian(data)

	if err := validateCreateUjian(data); err != nil {
		logger.Error(ctx, "failed create ujian", "layer", "core.service", "op", "ujian.create", "err", err)
		return err
	}

	if err := r.ujianRepo.CreateUjian(ctx, data); err != nil {
		logger.Error(ctx, "failed create ujian", "layer", "core.service", "op", "ujian.create", "err", err)
		return err
	}

	return nil
}

type CreatePesertaUjianService struct {
	pesertaRepo PesertaUjianRepository
}

func NewCreatePesertaUjianService(pesertaRepo PesertaUjianRepository) *CreatePesertaUjianService {
	return &CreatePesertaUjianService{
		pesertaRepo: pesertaRepo,
	}
}

func (r *CreatePesertaUjianService) CreatePesertaUjianService(ctx context.Context, data ujian.PesertaUjian) (ujian.ID, error) {
	logger := corelog.FromContext(ctx)

	if err := validateCreatePesertaUjian(data); err != nil {
		logger.Error(ctx, "failed create peserta ujian", "layer", "core.service", "op", "ujian.create_peserta", "err", err)
		return 0, err
	}

	id, err := r.pesertaRepo.CreatePesertaUjian(ctx, data)
	if err != nil {
		logger.Error(ctx, "failed create peserta ujian", "layer", "core.service", "op", "ujian.create_peserta", "err", err)
		return 0, err
	}

	return id, nil
}

type CreateJawabanUjianSiswaService struct {
	jawabanRepo JawabanUjianRepository
}

func NewCreateJawabanUjianSiswaService(jawabanRepo JawabanUjianRepository) *CreateJawabanUjianSiswaService {
	return &CreateJawabanUjianSiswaService{
		jawabanRepo: jawabanRepo,
	}
}

func (r *CreateJawabanUjianSiswaService) CreateJawabanUjianSiswaService(ctx context.Context, data ujian.JawabanUjianSiswa) (ujian.ID, error) {
	logger := corelog.FromContext(ctx)

	data = sanitizeCreateJawabanUjianSiswa(data)

	if err := validateCreateJawabanUjianSiswa(data); err != nil {
		logger.Error(ctx, "failed create jawaban ujian siswa", "layer", "core.service", "op", "ujian.create_jawaban", "err", err)
		return 0, err
	}

	id, err := r.jawabanRepo.CreateJawabanUjianSiswa(ctx, data)
	if err != nil {
		logger.Error(ctx, "failed create jawaban ujian siswa", "layer", "core.service", "op", "ujian.create_jawaban", "err", err)
		return 0, err
	}

	return id, nil
}
