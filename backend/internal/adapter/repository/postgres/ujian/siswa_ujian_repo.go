package ujianrepo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type SiswaUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewSiswaUjianRepo(q pg.Executor, logger corelog.Logger) *SiswaUjianRepo {
	return &SiswaUjianRepo{q: q, logger: logger}
}

func (r *SiswaUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *SiswaUjianRepo) CheckValidSiswaInPesertaUjianById(ctx context.Context, idSiswa int, idJadwalUjian int) (bool, int, error) {
	const query = `
		SELECT
			pu.id_peserta_ujian
		FROM peserta_ujian pu
		WHERE pu.id_siswa = $1
			AND pu.id_jadwal_ujian = $2
		LIMIT 1
	`

	var idPesertaUjian int

	err := r.q.QueryRow(ctx, query, idSiswa, idJadwalUjian).Scan(&idPesertaUjian)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, 0, nil
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed checking siswa in peserta ujian", "layer", "repo.db", "op", "siswa_ujian.check_peserta", "err", err)
		return false, 0, err
	}

	return true, idPesertaUjian, nil
}

func (r *SiswaUjianRepo) CheckTokenUjian(ctx context.Context, token string) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM jadwal_ujian
			WHERE UPPER(token) = UPPER($1)
		)
	`

	var exists bool
	if err := r.q.QueryRow(ctx, query, token).Scan(&exists); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed checking token ujian", "layer", "repo.db", "op", "siswa_ujian.check_token", "err", err)
		return false, err
	}

	return exists, nil
}

func (r *SiswaUjianRepo) GetDeadlineUjian(ctx context.Context, idJadwalUjian int) (time.Time, error) {
	const query = `
		SELECT waktu_selesai
		FROM jadwal_ujian
		WHERE id_jadwal_ujian = $1
	`

	var deadline time.Time
	err := r.q.QueryRow(ctx, query, idJadwalUjian).Scan(&deadline)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get deadline ujian", "layer", "repo.db", "op", "siswa_ujian.get_deadline", "err", err)
		return time.Time{}, err
	}

	return deadline, nil
}
