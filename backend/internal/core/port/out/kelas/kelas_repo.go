package kelas_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type KelasRepository interface {
	GetKelas(ctx context.Context, filter query.ListKelasFilter)([]kelas.FullKelasData, error)
}