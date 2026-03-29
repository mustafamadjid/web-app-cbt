package ujiansiswarepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type UjianSiswaRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewUjianSiswaRepo(q pg.Executor, logger corelog.Logger) *UjianSiswaRepo {
	return &UjianSiswaRepo{q: q, logger: logger}
}

func (r *UjianSiswaRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *UjianSiswaRepo) ListUjianSiswa(ctx context.Context, idSiswa int, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
	queryText, args := r.buildListUjianSiswaQuery(idSiswa, filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get list ujian siswa", "layer", "repo.db", "op", "siswa_ujian.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanListUjianSiswaRows(ctx, "siswa_ujian.list", rows)
}

func (r *UjianSiswaRepo) GetWaktuSelesaiUjian(ctx context.Context, idJadwalUjian int) (time.Time, error) {
	const queryText = `
		SELECT waktu_selesai
		FROM jadwal_ujian
		WHERE id_jadwal_ujian = $1
	`

	var waktuSelesai time.Time
	err := r.q.QueryRow(ctx, queryText, idJadwalUjian).Scan(&waktuSelesai)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get waktu selesai ujian", "layer", "repo.db", "op", "siswa_ujian.get_waktu_selesai", "err", err)
		return time.Time{}, err
	}

	return waktuSelesai, nil
}

func (r *UjianSiswaRepo) GetActiveUjianAttemptBySiswa(ctx context.Context, idSiswa int, idJadwalUjian int) (ujian.AttemptUjian, error) {
	const queryText = `
		WITH peserta AS (
			SELECT
				pu.id_peserta_ujian
			FROM peserta_ujian pu
			WHERE pu.id_siswa = $1
				AND pu.id_jadwal_ujian = $2
			LIMIT 1
		)
		SELECT
			au.id_attempt,
			au.id_peserta_ujian,
			au.status_attempt,
			au.waktu_mulai,
			au.waktu_submit,
			au.deadline_at
		FROM attempt_ujian au
		INNER JOIN peserta p ON p.id_peserta_ujian = au.id_peserta_ujian
		WHERE au.status_attempt = $3
		LIMIT 1
	`

	item, err := scanAttemptUjianRow(r.q.QueryRow(ctx, queryText, idSiswa, idJadwalUjian, ujian.ATTEMPT_IN_PROGRESS))
	if errors.Is(err, pgx.ErrNoRows) {
		return ujian.AttemptUjian{}, sql.ErrNoRows
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get active attempt ujian by siswa", "layer", "repo.db", "op", "attempt_ujian.get_active_by_siswa", "err", err)
		return ujian.AttemptUjian{}, err
	}

	return item, nil
}
