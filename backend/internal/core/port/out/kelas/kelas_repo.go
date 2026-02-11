package kelas_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type KelasRepository interface {
	GetKelas(ctx context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error)
	GetKelasById(ctx context.Context, idTingkatKelas int, idNamaKelas int) (kelas.KelasData, error)

	CreateTingkatKelas(ctx context.Context, tingkatKelas int) error
	CreateNamaKelas(ctx context.Context, namaKelas kelas.NamaKelas) error

	UpdateNamaKelas(ctx context.Context, idNamaKelas int, dataUpdate updatepatch.NamaKelasPatch) error

	ExistTingkatKelas(ctx context.Context, tingkatKelas int) (bool, error)
	ExistNamaKelas(ctx context.Context, namaKelas string) (bool, error)
}
