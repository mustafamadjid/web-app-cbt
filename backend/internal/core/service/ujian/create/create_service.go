package ujian_service

import (
	"context"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type CreateUjianService struct {
	repo ujian_repo.UjianRepository
}

func NewCreateUjianService(repo ujian_repo.UjianRepository) *CreateUjianService {
	return &CreateUjianService{
		repo: repo,
	}
}

func (r *CreateUjianService) CreateUjianService(ctx context.Context, data ujian.PenjadwalanUjian) error {
	logger := corelog.FromContext(ctx)

	data = sanitizeCreateUjian(data)

	if err := validateCreateUjian(data); err != nil {
		logger.Error(ctx, "failed create ujian", "layer", "core.service", "op", "ujian.create", "err", err)
		return err
	}

	if err := r.repo.CreateUjian(ctx, data); err != nil {
		logger.Error(ctx, "failed create ujian", "layer", "core.service", "op", "ujian.create", "err", err)
		return err
	}

	return nil
}

func (r *CreateUjianService) CreatePesertaUjianService(ctx context.Context, data ujian.PesertaUjian) (ujian.ID, error) {
	logger := corelog.FromContext(ctx)

	if err := validateCreatePesertaUjian(data); err != nil {
		logger.Error(ctx, "failed create peserta ujian", "layer", "core.service", "op", "ujian.create_peserta", "err", err)
		return 0, err
	}

	id, err := r.repo.CreatePesertaUjian(ctx, data)
	if err != nil {
		logger.Error(ctx, "failed create peserta ujian", "layer", "core.service", "op", "ujian.create_peserta", "err", err)
		return 0, err
	}

	return id, nil
}

func (r *CreateUjianService) CreateJawabanUjianSiswaService(ctx context.Context, data ujian.JawabanUjianSiswa) (ujian.ID, error) {
	logger := corelog.FromContext(ctx)

	data = sanitizeCreateJawabanUjianSiswa(data)

	if err := validateCreateJawabanUjianSiswa(data); err != nil {
		logger.Error(ctx, "failed create jawaban ujian siswa", "layer", "core.service", "op", "ujian.create_jawaban", "err", err)
		return 0, err
	}

	id, err := r.repo.CreateJawabanUjianSiswa(ctx, data)
	if err != nil {
		logger.Error(ctx, "failed create jawaban ujian siswa", "layer", "core.service", "op", "ujian.create_jawaban", "err", err)
		return 0, err
	}

	return id, nil
}
