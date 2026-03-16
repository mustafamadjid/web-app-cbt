package active_attempt_repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type ActiveAttemptRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewActiveAttemptRepo(q pg.Executor, logger corelog.Logger) *ActiveAttemptRepo {
	return &ActiveAttemptRepo{q: q, logger: logger}
}

func (r *ActiveAttemptRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *ActiveAttemptRepo) GetActiveUjianAttemptBySiswa(ctx context.Context, idSiswa int, idJadwalUjian int) (ujian.AttemptUjian, error) {
	const query = `
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

	var (
		item        ujian.AttemptUjian
		status      string
		waktuMulai  sql.NullTime
		waktuSubmit sql.NullTime
		deadlineAt  sql.NullTime
	)

	err := r.q.QueryRow(ctx, query, idSiswa, idJadwalUjian, ujian.ATTEMPT_IN_PROGRESS).Scan(
		&item.IdAttempt,
		&item.IdPesertaUjian,
		&status,
		&waktuMulai,
		&waktuSubmit,
		&deadlineAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ujian.AttemptUjian{}, sql.ErrNoRows
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get active attempt ujian by siswa", "layer", "repo.db", "op", "attempt_ujian.get_active_by_siswa", "err", err)
		return ujian.AttemptUjian{}, err
	}

	item.StatusAttempt = ujian.StatusAttempt(status)
	item.WaktuMulai = nullTimeToPtr(waktuMulai)
	item.WaktuSubmit = nullTimeToPtr(waktuSubmit)
	item.DeadlineAt = nullTimeToPtr(deadlineAt)

	return item, nil
}

func nullTimeToPtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	v := value.Time
	return &v
}
