package ujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type CreateUjianService struct {
	ujianRepo ujian_repo.UjianRepository
}

func NewCreateUjianService(ujianRepo ujian_repo.UjianRepository) *CreateUjianService {
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
