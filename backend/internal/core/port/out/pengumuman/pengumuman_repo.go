package pengumuman_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type PengumumanRepo interface {
	GetPengumumanActive(ctx context.Context)([]pengumuman.Pengumuman, error)
	GetPengumumanNonActive(ctx context.Context)([]pengumuman.Pengumuman, error)
	GetPengumumanIncoming(ctx context.Context)([]pengumuman.Pengumuman, error)
	GetPengumumanById(ctx context.Context, idPengumuman pengumuman.ID)(pengumuman.Pengumuman, error)
	
	CreatePengumuman(ctx context.Context, pengumuman pengumuman.Pengumuman)error
	UpdatePengumuman(ctx context.Context, idPengumuman pengumuman.ID, dataUpdate updatepatch.PengumumanUpdatePatch)error
	DeletePengumuman(ctx context.Context, idPengumuman pengumuman.ID)error
}