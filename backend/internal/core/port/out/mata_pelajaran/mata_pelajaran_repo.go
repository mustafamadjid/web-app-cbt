package matapelajaran_repo

import (
	"context"

	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type MataPelajaranRepository interface {
	GetMapel(ctx context.Context, filter query.ListMapelFilter) ([]matapelajaran.MataPelajaran, error)
	GetMapelById(ctx context.Context, idMapel int) (matapelajaran.MataPelajaran, error)

	CreateMapel(ctx context.Context, mapel matapelajaran.MataPelajaran) error

	UpdateMapel(ctx context.Context,idMapel int, mapel updatepatch.UpdateMapelPatch) error

	DeleteMapel(ctx context.Context, idMapel int) error

	ExistKodeMapel(ctx context.Context, kodeMapel string)(bool,error)
}