package gradingrepo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type GradingRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

type statistikSoalUpsertItem struct {
	IDSoal             ujian.ID `json:"id_soal"`
	IDUjian            ujian.ID `json:"id_ujian"`
	JumlahJawabanBenar int      `json:"jumlah_jawaban_benar"`
	JumlahJawabanSalah int      `json:"jumlah_jawaban_salah"`
}

func NewGradingRepo(q pg.Executor, logger corelog.Logger) *GradingRepo {
	return &GradingRepo{q: q, logger: logger}
}

func (r *GradingRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *GradingRepo) UpsertNilaiToHasilUjian(ctx context.Context, totalNilai float64, hasilUjian ujian.HasilUjian) error {
	const query = `
		INSERT INTO hasil_ujian (
			id_attempt,
			graded_by,
			nilai_akhir,
			passed,
			essay_graded,
			graded_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			COALESCE($5, FALSE),
			COALESCE($6, NOW())
		)
		ON CONFLICT (id_attempt)
		DO UPDATE SET
			graded_by = COALESCE($2, hasil_ujian.graded_by),
			nilai_akhir = $3,
			passed = COALESCE($4, hasil_ujian.passed),
			essay_graded = COALESCE($5, hasil_ujian.essay_graded),
			graded_at = COALESCE($6, NOW())
	`

	_, err := r.q.Exec(
		ctx,
		query,
		hasilUjian.IdAttempt,
		hasilUjian.GradedBy,
		totalNilai,
		hasilUjian.Passed,
		hasilUjian.EssayGraded,
		hasilUjian.GradedAt,
	)
	if err != nil {
		if mappedErr := mapGradingConstraintError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed upsert nilai hasil ujian", "layer", "repo.db", "op", "ujian.grading.upsert_nilai", "attempt_id", hasilUjian.IdAttempt, "err", err)
		return err
	}

	return nil
}

func (r *GradingRepo) UpsertJawabanBenarToStatistikSoal(ctx context.Context, soalBenar []ujian.StatistikSoal) error {
	payload := buildStatistikSoalPayload(soalBenar, true)
	return r.upsertStatistikSoal(ctx, payload)
}

func (r *GradingRepo) UpsertJawabanSalahToStatistikSoal(ctx context.Context, soalSalah []ujian.StatistikSoal) error {
	payload := buildStatistikSoalPayload(soalSalah, false)
	return r.upsertStatistikSoal(ctx, payload)
}

func (r *GradingRepo) upsertStatistikSoal(ctx context.Context, payload []statistikSoalUpsertItem) error {
	if len(payload) == 0 {
		return nil
	}

	const query = `
		WITH aggregated AS (
			SELECT
				x.id_soal,
				x.id_ujian,
				SUM(x.jumlah_jawaban_benar)::integer AS jumlah_jawaban_benar,
				SUM(x.jumlah_jawaban_salah)::integer AS jumlah_jawaban_salah
			FROM jsonb_to_recordset($1::jsonb) AS x(
				id_soal bigint,
				id_ujian bigint,
				jumlah_jawaban_benar integer,
				jumlah_jawaban_salah integer
			)
			GROUP BY x.id_soal, x.id_ujian
		)
		INSERT INTO statistik_soal (
			id_soal,
			id_ujian,
			jumlah_jawaban_benar,
			jumlah_jawaban_salah
		)
		SELECT
			id_soal,
			id_ujian,
			jumlah_jawaban_benar,
			jumlah_jawaban_salah
		FROM aggregated
		ON CONFLICT (id_ujian, id_soal)
		DO UPDATE SET
			jumlah_jawaban_benar = statistik_soal.jumlah_jawaban_benar + EXCLUDED.jumlah_jawaban_benar,
			jumlah_jawaban_salah = statistik_soal.jumlah_jawaban_salah + EXCLUDED.jumlah_jawaban_salah,
			updated_at = NOW()
	`

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed marshal statistik soal payload", "layer", "repo.db", "op", "ujian.grading.upsert_statistik.marshal", "err", err)
		return err
	}

	_, err = r.q.Exec(ctx, query, payloadJSON)
	if err != nil {
		if mappedErr := mapGradingConstraintError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed upsert statistik soal", "layer", "repo.db", "op", "ujian.grading.upsert_statistik.exec", "err", err)
		return err
	}

	return nil
}

func mapGradingConstraintError(err error) error {
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
