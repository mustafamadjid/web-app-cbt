package pengumuman_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
)

type PengumumanRepo interface {
	GetPengumuman(ctx context.Context)([]pengumuman.Pengumuman, error)
	GetPengumumanById(ctx context.Context, idPengumuman int)(pengumuman.Pengumuman, error)
	
	CreatePengumuman(ctx context.Context, pengumuman pengumuman.Pengumuman)(int, error)
	UpdatePengumuman(ctx context.Context, idPengumuman int, dataUpdate pengumuman.Pengumuman)(int, error)
	DeletePengumuman(ctx context.Context, idPengumuman int)(int, error)
}