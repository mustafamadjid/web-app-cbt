package jawabanujian_repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type JawabanUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewJawabanUjianRepo(q pg.Executor, logger corelog.Logger) *JawabanUjianRepo {
	return &JawabanUjianRepo{q: q, logger: logger}
}

func (r *JawabanUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *JawabanUjianRepo) SaveJawabanUjian(ctx context.Context, idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error {
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
	// Marshalling
	payloadJson, err := json.Marshal(jawaban)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed marshal jawaban ujian", "layer", "adapter.repository", "op", "ujian.jawaban.save", "err", err)
		return err
	}

	_, err = r.q.Exec(ctx, query, idAttempt, payloadJson)
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

	var jawaban []ujian.JawabanUjian
	for rows.Next() {
		var (
			itemJawaban  ujian.JawabanUjian
			IdPilihan    sql.NullInt64
			JawabanEssay sql.NullString
			WaktuJawab   sql.NullTime
		)

		if err := rows.Scan(
			&itemJawaban.IdJawaban,
			&itemJawaban.IdSoal,
			&IdPilihan,
			&JawabanEssay,
			&WaktuJawab,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scan jawaban ujian by attempt id", "layer", "adapter.repository", "op", "ujian.jawaban.get_by_attempt_id", "err", err)
			return nil, err
		}

		itemJawaban.IdPilihan = nullInt64ToUjianIDPtr(IdPilihan)
		itemJawaban.JawabanEssay = nullStringToPtr(JawabanEssay)

		if WaktuJawab.Valid {
			itemJawaban.WaktuJawab = &WaktuJawab.Time
		}

		jawaban = append(jawaban, itemJawaban)
	}
	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get jawaban ujian by attempt id", "layer", "adapter.repository", "op", "ujian.jawaban.get_by_attempt_id", "err", err)
		return nil, err
	}

	return jawaban, nil
}

func nullInt64ToUjianIDPtr(v sql.NullInt64) *ujian.ID {
	if !v.Valid {
		return nil
	}
	id := ujian.ID(v.Int64)
	return &id
}

func nullStringToPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
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
