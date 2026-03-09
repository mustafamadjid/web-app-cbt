package ujian_service

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type updateUjianRepository interface {
	UpdateUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error
}

type UjianRepository struct {
	repo updateUjianRepository
}

func NewUjianRepository(repo updateUjianRepository) UjianRepository {
	return UjianRepository{repo: repo}
}

func (r UjianRepository) UpdateUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error {
	return r.repo.UpdateUjian(ctx, id, payload)
}
