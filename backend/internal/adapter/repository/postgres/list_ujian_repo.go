package postgres

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
)

type ListUjianRepo struct {
	q Executor
	logger corelog.Logger
}

func NewListUjianRepo(q Executor, logger corelog.Logger) *ListUjianRepo {
	return &ListUjianRepo{q:q, logger: logger}
}

func (r *ListUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func(r *ListUjianRepo)GetAllUjian(ctx context.Context, filter query.ListRuangUjianFilter) ([]ujian.Ujian, error){
	baseQuery := `
		SELECT
			u.id_ujian,
			u.nama_ujian,
			u.id_kelas,
			u.id_nama_kelas,
			ju.id_jadwal_ujian,
			ju.tanggal_ujian,
			ju.waktu_mulai,
			ju.waktu_selesai,
			ju.status_ujian
		FROM ujian u
		JOIN jadwal_ujian ju
		ON ju.id_ujian = u.id_ujian;
	`
}