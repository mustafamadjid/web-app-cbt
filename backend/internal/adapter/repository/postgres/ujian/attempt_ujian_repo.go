package ujianrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func (r *UjianRepo) GetAttemptUjianById(ctx context.Context, idAttempt ujian.ID) (ujian.AttemptUjian, error) {
	const query = `
		SELECT
			id_attempt,
			id_peserta_ujian,
			status_attempt,
			waktu_mulai,
			waktu_submit,
			deadline_at
		FROM attempt_ujian
		WHERE id_attempt = $1
	`

	var (
		item        ujian.AttemptUjian
		status      string
		waktuMulai  sql.NullTime
		waktuSubmit sql.NullTime
		deadlineAt  sql.NullTime
	)

	err := r.q.QueryRow(ctx, query, idAttempt).Scan(
		&item.IdAttempt,
		&item.IdPesertaUjian,
		&status,
		&waktuMulai,
		&waktuSubmit,
		&deadlineAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return ujian.AttemptUjian{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get attempt ujian by id", "layer", "repo.db", "op", "attempt_ujian.get_by_id", "err", err)
		return ujian.AttemptUjian{}, err
	}

	item.StatusAttempt = ujian.StatusAttempt(status)
	item.WaktuMulai = nullTimeToPtr(waktuMulai)
	item.WaktuSubmit = nullTimeToPtr(waktuSubmit)
	item.DeadlineAt = nullTimeToPtr(deadlineAt)

	return item, nil
}

func (r *UjianRepo) CreateAttemptUjian(ctx context.Context, data ujian.AttemptUjian) error {
	const query = `
		INSERT INTO attempt_ujian (
			id_peserta_ujian,
			status_attempt,
			waktu_mulai,
			waktu_submit,
			deadline_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	status := data.StatusAttempt
	if status == "" {
		status = ujian.ATTEMPT_IN_PROGRESS
	}

	_, err := r.q.Exec(
		ctx,
		query,
		data.IdPesertaUjian,
		status,
		timePtrToDB(data.WaktuMulai),
		timePtrToDB(data.WaktuSubmit),
		timePtrToDB(data.DeadlineAt),
	)
	if err != nil {
		if mappedErr := mapAttemptUniqueViolation(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed create attempt ujian", "layer", "repo.db", "op", "attempt_ujian.create", "err", err)
		return err
	}

	return nil
}

func (r *UjianRepo) UpdateAttemptUjian(ctx context.Context, idAttempt ujian.ID, data updatepatch.UpdateAttemptUjianPatch) error {
	const query = `
		UPDATE attempt_ujian
		SET
			id_peserta_ujian = COALESCE($1, id_peserta_ujian),
			status_attempt = COALESCE($2, status_attempt),
			waktu_mulai = COALESCE($3, waktu_mulai),
			waktu_submit = COALESCE($4, waktu_submit),
			deadline_at = COALESCE($5, deadline_at),
			updated_at = COALESCE($6, NOW())
		WHERE id_attempt = $7
	`

	tag, err := r.q.Exec(
		ctx,
		query,
		data.IdPesertaUjian,
		data.StatusAttempt,
		timePtrToDB(data.WaktuMulai),
		timePtrToDB(data.WaktuSubmit),
		timePtrToDB(data.DeadlineAt),
		timePtrToDB(data.UpdatedAt),
		idAttempt,
	)
	if err != nil {
		if mappedErr := mapAttemptUniqueViolation(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed update attempt ujian", "layer", "repo.db", "op", "attempt_ujian.update", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *UjianRepo) DeleteAttemptUjian(ctx context.Context, idAttempt ujian.ID) error {
	const query = `
		DELETE FROM attempt_ujian
		WHERE id_attempt = $1
	`

	tag, err := r.q.Exec(ctx, query, idAttempt)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed delete attempt ujian", "layer", "repo.db", "op", "attempt_ujian.delete", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func mapAttemptUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}

	return coreerror.ErrConflict
}

func nullTimeToPtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	v := value.Time
	return &v
}

func timePtrToDB(value *time.Time) any {
	if value == nil {
		return nil
	}

	return *value
}
