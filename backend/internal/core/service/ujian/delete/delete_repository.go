package ujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type deleteUjianRepository interface {
	DeleteUjian(ctx context.Context, id ujian.ID) error
}

type UjianRepository struct {
	repo deleteUjianRepository
}

func NewUjianRepository(repo deleteUjianRepository) UjianRepository {
	return UjianRepository{repo: repo}
}

func (r UjianRepository) DeleteUjian(ctx context.Context, id ujian.ID) error {
	return r.repo.DeleteUjian(ctx, id)
}
