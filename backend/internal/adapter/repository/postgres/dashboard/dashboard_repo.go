package dashboard_repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/dashboard"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type DashboardStatistikRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewDashboardStatistikRepo(q pg.Executor, logger corelog.Logger) *DashboardStatistikRepo {
	return &DashboardStatistikRepo{q: q, logger: logger}
}

func (r *DashboardStatistikRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *DashboardStatistikRepo) GetDashboardStatistik(ctx context.Context) (dashboard.DashboardStatistik, error) {
	const query = `
		SELECT
			(SELECT COUNT(*) FROM profil_siswa) AS total_siswa,
			(SELECT COUNT(*) FROM profil_guru) AS total_guru,
			(
				SELECT COUNT(*)
				FROM jadwal_ujian ju
				WHERE ju.waktu_selesai < NOW()
			) AS total_ujian_terlaksana,

			(SELECT COUNT(*) FROM bank_soal) AS total_bank_soal,

			(
				SELECT COUNT(*)
				FROM mata_pelajaran mp
			) AS total_mapel_aktif
	`

	var statistik dashboard.DashboardStatistik
	err := r.q.QueryRow(ctx, query).Scan(
		&statistik.TotalSiswa,
		&statistik.TotalGuru,
		&statistik.TotalUjianTerlaksana,
		&statistik.TotalBankSoal,
		&statistik.TotalMapelAktif,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dashboard.DashboardStatistik{}, nil
		}
		r.loggerFor(ctx).Error(ctx, "failed getting dashboard statistik", "layer", "repo.db", "op", "dashboard.get_statistik", "err", err)
		return dashboard.DashboardStatistik{}, err
	}

	return statistik, nil
}
