package postgres

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type KelasRepo struct {
	q  Executor
	logger corelog.Logger
}

func NewKelasRepo(q Executor, logger corelog.Logger) *KelasRepo {
	return &KelasRepo{q: q, logger: logger}
}

func(r *KelasRepo)loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func(r *KelasRepo)GetKelas(ctx context.Context, filter query.ListKelasFilter)([]kelas.FullKelasData, error) {
	baseQuery := `
	SELECT 
	tk.id_kelas,
	tk.tingkat_kelas,
	k.id_nama_kelas,
	k.nama_kelas
	FROM kelas tk
	JOIN nama_kelas k ON tk.id_kelas = k.id_kelas
	`
}