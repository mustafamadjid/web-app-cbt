package kelas_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type KelasRepository interface {
	GetKelas(ctx context.Context, filter query.ListKelasFilter)([]kelas.FullKelasData, error)
	CreateTingkatKelas(ctx context.Context, tingkatKelas kelas.TingkatKelas)(kelas.ID,error)
	CreateNamaKelas(ctx context.Context, namaKelas kelas.NamaKelas)(kelas.ID,error)
}