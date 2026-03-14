package siswaujianrepo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

func (r *SiswaUjianRepo) GetWaktuSelesaiUjian(ctx context.Context, idJadwalUjian int) (time.Time, error) {
	const query = `
		SELECT waktu_selesai
		FROM jadwal_ujian
		WHERE id_jadwal_ujian = $1
	`

	var waktuSelesai time.Time
	err := r.q.QueryRow(ctx, query, idJadwalUjian).Scan(&waktuSelesai)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get waktu selesai ujian", "layer", "repo.db", "op", "siswa_ujian.get_waktu_selesai", "err", err)
		return time.Time{}, err
	}

	return waktuSelesai, nil
}
