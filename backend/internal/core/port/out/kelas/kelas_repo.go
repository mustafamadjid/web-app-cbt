package kelas_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type KelasRepository interface {
	GetKelas(ctx context.Context, filter query.ListKelasFilter)([]kelas.FullKelasData, error)
	CreateTingkatKelas(ctx context.Context, tingkatKelas int)error
	CreateNamaKelas(ctx context.Context, namaKelas kelas.NamaKelas)error

	ExistTingkatKelas(ctx context.Context, tingkatKelas int )(bool, error)
	ExistNamaKelas(ctx context.Context,namaKelas string)(bool, error)
}