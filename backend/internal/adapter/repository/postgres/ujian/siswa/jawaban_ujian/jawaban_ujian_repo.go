package jawabanujian_repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)



type JawabanUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
	pool   *pgxpool.Pool
}

func NewJawabanUjianRepo(q pg.Executor, logger corelog.Logger, pool *pgxpool.Pool) *JawabanUjianRepo {
	return &JawabanUjianRepo{q: q, logger: logger, pool: pool}
}

func (r *JawabanUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *JawabanUjianRepo) SaveJawabanUjian(ctx context.Context, idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error {
	upsertItems, clearSoalIDs := splitSaveJawabanItems(jawaban)
	if len(upsertItems) == 0 && len(clearSoalIDs) == 0 {
		return nil
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed begin tx save jawaban ujian", "layer", "adapter.repository", "op", "ujian.jawaban.save.begin_tx", "err", err)
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	if len(clearSoalIDs) > 0 {
		if _, err := tx.Exec(
			ctx,
			`DELETE FROM jawaban_ujian_siswa
			WHERE id_attempt = $1 AND id_soal = ANY($2)`,
			idAttempt,
			clearSoalIDs,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed clear jawaban ujian", "layer", "adapter.repository", "op", "ujian.jawaban.clear", "err", err)
			return err
		}
	}

	if len(upsertItems) > 0 {
		if err := r.upsertJawabanUjian(ctx, tx, idAttempt, upsertItems); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit save jawaban ujian", "layer", "adapter.repository", "op", "ujian.jawaban.save.commit", "err", err)
		return err
	}

	committed = true
	return nil
}

func (r *JawabanUjianRepo) upsertJawabanUjian(ctx context.Context, q pg.Executor, idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error {
	query := `
		INSERT INTO jawaban_ujian_siswa (
			id_attempt,
			id_soal,
			id_pilihan,
			jawaban_essay,
			waktu_jawab
			)
			SELECT
			$1,
			x.id_soal,
			x.id_pilihan,
			x.jawaban_essay,
			x.waktu_jawab
		FROM jsonb_to_recordset($2::jsonb) AS x(
			id_soal bigint,
			id_pilihan bigint,
			jawaban_essay text,
			waktu_jawab timestamptz
			)
		ON CONFLICT (id_attempt, id_soal)
		DO UPDATE SET
			id_pilihan = EXCLUDED.id_pilihan,
			jawaban_essay = EXCLUDED.jawaban_essay,
			waktu_jawab = EXCLUDED.waktu_jawab
		`

	payloadJson, err := json.Marshal(jawaban)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed marshal jawaban ujian", "layer", "adapter.repository", "op", "ujian.jawaban.save", "err", err)
		return err
	}

	_, err = q.Exec(ctx, query, idAttempt, payloadJson)
	if err != nil {
		if mappedErr := mapJawabanConstraintError(err); mappedErr != nil {
			return mappedErr
		}
		r.loggerFor(ctx).Error(ctx, "failed save jawaban ujian", "layer", "adapter.repository", "op", "ujian.jawaban.save", "err", err)
		return err
	}
	return nil
}

func (r *JawabanUjianRepo) GetJawabanUjianByAttemptId(ctx context.Context, idAttempt ujian.ID) ([]ujian.JawabanUjian, error) {

	query := `
		SELECT
			id_jawaban,
			id_soal,
			id_pilihan,
			jawaban_essay,
			waktu_jawab
		FROM jawaban_ujian_siswa
		WHERE id_attempt = $1
	`
	rows, err := r.q.Query(ctx, query, idAttempt)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get jawaban ujian by attempt id", "layer", "adapter.repository", "op", "ujian.jawaban.get_by_attempt_id", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanJawabanUjianRows(ctx, "ujian.jawaban.get_by_attempt_id", rows)
}

func mapJawabanConstraintError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {
	case "23502", "23503", "23514":
		return coreerror.ErrInvalidInput
	case "23505":
		return coreerror.ErrConflict
	default:
		return nil
	}
}
