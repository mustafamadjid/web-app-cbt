package gradingrepo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type GradingRepo struct {
	q      pg.Executor
	logger corelog.Logger
	pool   *pgxpool.Pool
}

type statistikSoalUpsertItem struct {
	IDSoal             ujian.ID `json:"id_soal"`
	IDUjian            ujian.ID `json:"id_ujian"`
	JumlahJawabanBenar int      `json:"jumlah_jawaban_benar"`
	JumlahJawabanSalah int      `json:"jumlah_jawaban_salah"`
}

func NewGradingRepo(q pg.Executor, logger corelog.Logger) *GradingRepo {
	repo := &GradingRepo{q: q, logger: logger}
	if pool, ok := q.(*pgxpool.Pool); ok {
		repo.pool = pool
	}

	return repo
}

func (r *GradingRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *GradingRepo) ClaimQueuedJobs(ctx context.Context, limit int) ([]ujian.GradingJob, error) {
	const query = `
		UPDATE grading_jobs gj
		SET
			status = $1::varchar,
			error_message = NULL,
			error_code = NULL,
			locked_at = NOW(),
			updated_at = NOW()
		FROM (
			SELECT id_grading_jobs
			FROM grading_jobs
			WHERE status = $2
				AND available_at <= NOW()
			ORDER BY available_at ASC, id_grading_jobs ASC
			FOR UPDATE SKIP LOCKED
			LIMIT $3
		) claimed
		WHERE gj.id_grading_jobs = claimed.id_grading_jobs
		RETURNING
			gj.id_grading_jobs,
			gj.id_attempt,
			gj.status,
			gj.retry_count,
			gj.max_retries,
			COALESCE(gj.available_at::text, ''),
			COALESCE(gj.locked_at::text, ''),
			COALESCE(gj.error_code, ''),
			COALESCE(gj.error_message, '')
	`

	rows, err := r.q.Query(ctx, query, ujian.StatusProcessing, ujian.StatusQueued, limit)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed claiming queued grading jobs", "layer", "repo.db", "op", "ujian.grading.worker.claim_jobs", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanGradingJobRows(ctx, "ujian.grading.worker.claim_jobs", rows)
}

func (r *GradingRepo) UpdateStatusJob(ctx context.Context, jobID int, statusJob ujian.JobStatus, errorMsg string, errorCode string) error {
	const query = `
		UPDATE grading_jobs
		SET
			status = $1::varchar,
			error_message = CASE
				WHEN $1::varchar IN ('processing', 'done') THEN NULL
				ELSE NULLIF($2::varchar, '')
			END,
			error_code = CASE
				WHEN $1::varchar IN ('processing', 'done') THEN NULL
				ELSE NULLIF($3::varchar, '')
			END,
			locked_at = CASE
				WHEN $1::varchar = 'processing' THEN NOW()
				WHEN $1::varchar IN ('done', 'failed') THEN NULL
				ELSE locked_at
			END,
			updated_at = NOW()
		WHERE id_grading_jobs = $4
			AND (
				($1::varchar IN ('done', 'failed') AND status = 'processing')
				OR ($1::varchar NOT IN ('done', 'failed'))
			)
	`

	tag, err := r.q.Exec(ctx, query, statusJob, errorMsg, errorCode, jobID)
	if err != nil {
		if mappedErr := mapGradingConstraintError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed updating grading job status", "layer", "repo.db", "op", "ujian.grading.worker.update_status", "job_id", jobID, "status", statusJob, "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *GradingRepo) UpsertNilaiToHasilUjian(ctx context.Context, totalNilai float64, hasilUjian ujian.HasilUjian) error {
	const query = `
		WITH data_ujian AS (
			SELECT ju.id_jadwal_ujian
			FROM attempt_ujian au
			INNER JOIN peserta_ujian pu 
				ON pu.id_peserta_ujian = au.id_peserta_ujian
			INNER JOIN jadwal_ujian ju 
				ON ju.id_jadwal_ujian = pu.id_jadwal_ujian
			WHERE au.id_attempt = $1
		), essay_score AS (
			SELECT
				COALESCE(SUM(CASE WHEN jus.essay_is_benar = true THEN s.bobot_soal ELSE 0 END), 0)::decimal(5,2) AS tambahan_nilai
			FROM jawaban_ujian_siswa jus
			INNER JOIN isi_soal s
				ON s.id_soal = jus.id_soal
			WHERE jus.id_attempt = $1
		), hasil_upsert AS (
			INSERT INTO hasil_ujian (
				id_attempt,
				graded_by,
				nilai_akhir,
				passed,
				essay_graded,
				graded_at,
				id_jadwal_ujian
			)
			VALUES (
				$1,
				$2,
				$3 + COALESCE((SELECT tambahan_nilai FROM essay_score), 0),
				$4,
				COALESCE($5, FALSE),
				COALESCE($6, NOW()),
				(SELECT id_jadwal_ujian FROM data_ujian)
			)
			ON CONFLICT (id_attempt)
			DO UPDATE SET
				graded_by = COALESCE($2, hasil_ujian.graded_by),
				nilai_akhir = $3 + COALESCE((SELECT tambahan_nilai FROM essay_score), 0),
				passed = COALESCE($4, hasil_ujian.passed),
				essay_graded = COALESCE($5, hasil_ujian.essay_graded),
				graded_at = COALESCE($6, hasil_ujian.graded_at),
				id_jadwal_ujian = COALESCE(
					(SELECT id_jadwal_ujian FROM data_ujian),
					hasil_ujian.id_jadwal_ujian
				)
			RETURNING
				id_attempt,
				id_jadwal_ujian,
				nilai_akhir
		), statistik_ujian_source AS (
			SELECT
				hu.id_jadwal_ujian,
				hu.nilai_akhir
			FROM hasil_upsert hu
			WHERE hu.nilai_akhir IS NOT NULL
			UNION ALL
			SELECT
				existing.id_jadwal_ujian,
				existing.nilai_akhir
			FROM hasil_ujian existing
			INNER JOIN hasil_upsert hu
				ON hu.id_jadwal_ujian = existing.id_jadwal_ujian
			WHERE existing.id_attempt <> hu.id_attempt
				AND existing.nilai_akhir IS NOT NULL
		), statistik_ujian_agg AS (
			SELECT
				sus.id_jadwal_ujian,
				MAX(sus.nilai_akhir)::decimal(5,2) AS nilai_tertinggi,
				MIN(sus.nilai_akhir)::decimal(5,2) AS nilai_terendah,
				ROUND(AVG(sus.nilai_akhir)::numeric, 2)::decimal(5,2) AS nilai_rata_rata,
				COUNT(*)::integer AS total_peserta_ujian
			FROM statistik_ujian_source sus
			GROUP BY sus.id_jadwal_ujian
		)
		INSERT INTO statistik_ujian (
			id_jadwal_ujian,
			nilai_tertinggi,
			nilai_terendah,
			nilai_rata_rata,
			total_peserta_ujian
		)
		SELECT
			sua.id_jadwal_ujian,
			sua.nilai_tertinggi,
			sua.nilai_terendah,
			sua.nilai_rata_rata,
			sua.total_peserta_ujian
		FROM statistik_ujian_agg sua
		ON CONFLICT (id_jadwal_ujian)
		DO UPDATE SET
			nilai_tertinggi = EXCLUDED.nilai_tertinggi,
			nilai_terendah = EXCLUDED.nilai_terendah,
			nilai_rata_rata = EXCLUDED.nilai_rata_rata,
			total_peserta_ujian = EXCLUDED.total_peserta_ujian,
			updated_at = NOW();
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
