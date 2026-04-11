package statistikujian_repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type StatistikUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewStatistikUjianRepo(q pg.Executor, logger corelog.Logger) *StatistikUjianRepo {
	return &StatistikUjianRepo{q: q, logger: logger}
}

func (r *StatistikUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *StatistikUjianRepo) GetStatistikUjianByIdJadwal(ctx context.Context, idJadwal ujian.ID) (ujian.StatistikUjian, error) {
	const query = `
		SELECT 
			id_statistik_ujian,
			id_jadwal_ujian,
			nilai_tertinggi,
			nilai_terendah,
			nilai_rata_rata,
			total_peserta_ujian
		FROM statistik_ujian
		WHERE id_jadwal_ujian = $1
	`

	var item ujian.StatistikUjian

	err := r.q.QueryRow(ctx, query, idJadwal).Scan(
		&item.IDStatistikUjian,
		&item.IDJadwalUjian,
		&item.NilaiTertinggi,
		&item.NilaiTerendah,
		&item.NilaiRataRata,
		&item.TotalPesertaUjian,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ujian.StatistikUjian{}, coreerror.ErrNotFound
		}
		if mappedErr := mapStatistikUjianConstraintError(err); mappedErr != nil {
			return ujian.StatistikUjian{}, mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed get statistik ujian by jadwal", "layer", "repo.db", "op", "statistik_ujian.get_by_jadwal", "jadwal_id", idJadwal, "err", err)
		return ujian.StatistikUjian{}, err
	}

	return item, nil
}

func mapStatistikUjianConstraintError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}

	return coreerror.ErrConflict
}
