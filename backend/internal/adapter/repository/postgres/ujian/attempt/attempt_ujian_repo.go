package attemptrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type AttemptUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
	pool   *pgxpool.Pool
}

func NewAttemptUjianRepo(q pg.Executor, logger corelog.Logger) *AttemptUjianRepo {
	repo := &AttemptUjianRepo{q: q, logger: logger}
	if pool, ok := q.(*pgxpool.Pool); ok {
		repo.pool = pool
	}

	return repo
}

func (r *AttemptUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *AttemptUjianRepo) GetAttemptUjianById(ctx context.Context, idAttempt ujian.ID) (ujian.AttemptUjian, error) {
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

	item, err := scanAttemptUjianRow(r.q.QueryRow(ctx, query, idAttempt))
	if errors.Is(err, pgx.ErrNoRows) {
		return ujian.AttemptUjian{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get attempt ujian by id", "layer", "repo.db", "op", "attempt_ujian.get_by_id", "err", err)
		return ujian.AttemptUjian{}, err
	}

	return item, nil
}

func (r *AttemptUjianRepo) CreateAttemptUjian(ctx context.Context, data ujian.AttemptUjian) error {
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

func (r *AttemptUjianRepo) UpdateAttemptUjian(ctx context.Context, idAttempt ujian.ID, data updatepatch.UpdateAttemptUjianPatch) error {
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

func (r *AttemptUjianRepo) DeleteAttemptUjian(ctx context.Context, idAttempt ujian.ID) error {
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

func (r *AttemptUjianRepo) SubmitAttemptUjian(ctx context.Context, idAttempt ujian.ID) error {
	if r.pool == nil {
		err := errors.New("attempt ujian repo requires pgx pool for submit transaction")
		r.loggerFor(ctx).Error(ctx, "failed submit attempt ujian", "layer", "repo.db", "op", "attempt_ujian.submit.precondition", "err", err)
		return err
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed begin tx submit attempt ujian", "layer", "repo.db", "op", "attempt_ujian.submit.begin_tx", "err", err)
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	const submitQuery = `
		UPDATE attempt_ujian
		SET
			status_attempt = $1,
			waktu_submit = NOW()
		WHERE id_attempt = $2 AND status_attempt = $3
	`

	tag, err := tx.Exec(ctx, submitQuery,
		ujian.ATTEMPT_SUBMITTED,
		idAttempt,
		ujian.ATTEMPT_IN_PROGRESS,
	)
	if err != nil {
		if mappedErr := mapSubmitUniqueViolation(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed update attempt ujian", "layer", "repo.db", "op", "attempt_ujian.submit.update", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound // atau ErrInvalidState, tergantung kebutuhan
	}

	const gradingJobsQuery = `
		INSERT INTO grading_jobs(id_attempt, status)
		VALUES ($1, 'queued')
		ON CONFLICT (id_attempt) DO NOTHING
	`

	_, err = tx.Exec(ctx, gradingJobsQuery, idAttempt)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed insert grading job", "layer", "repo.db", "op", "attempt_ujian.submit.insert_grading_job", "err", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit tx submit attempt ujian", "layer", "repo.db", "op", "attempt_ujian.submit.commit_tx", "err", err)
		return err
	}
	committed = true

	return nil
}
func mapAttemptUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}

	if pgErr.ConstraintName == "uq_attempt_active" {
		return coreerror.ErrSiswaHasActiveAttempt
	}

	return coreerror.ErrConflict
}

func mapSubmitUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}

	if pgErr.ConstraintName == "uq_attempt_ujian_one_submitted_per_peserta" {
		return coreerror.ErrDoubleSubmit
	}

	return coreerror.ErrConflict
}
