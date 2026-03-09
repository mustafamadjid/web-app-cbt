package ujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type createUjianRepository interface {
	CreateUjian(ctx context.Context, ujian ujian.PenjadwalanUjian) error
}

type UjianRepository struct {
	repo createUjianRepository
}

func NewUjianRepository(repo createUjianRepository) UjianRepository {
	return UjianRepository{repo: repo}
}

func (r UjianRepository) CreateUjian(ctx context.Context, data ujian.PenjadwalanUjian) error {
	return r.repo.CreateUjian(ctx, data)
}
