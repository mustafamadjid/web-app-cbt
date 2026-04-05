package gradingrepo

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func (r *GradingRepo) UpsertToStatistikUjian(ctx context.Context, idAttempt ujian.ID) error {
	const query = `
		WITH target_jadwal AS (
			SELECT
				pu.id_jadwal_ujian
			FROM attempt_ujian au
			JOIN peserta_ujian pu
				ON pu.id_peserta_ujian = au.id_peserta_ujian
			WHERE au.id_attempt = $1
		),
		agg_hasil_ujian AS (
			SELECT
				tj.id_jadwal_ujian,
				MAX(hu.nilai_akhir)::decimal(5,2) AS nilai_tertinggi,
				MIN(hu.nilai_akhir)::decimal(5,2) AS nilai_terendah,
				ROUND(AVG(hu.nilai_akhir)::numeric, 2)::decimal(5,2) AS nilai_rata_rata
			FROM target_jadwal tj
			JOIN hasil_ujian hu
				ON hu.id_jadwal_ujian = tj.id_jadwal_ujian
				AND hu.nilai_akhir IS NOT NULL
			GROUP BY tj.id_jadwal_ujian
		),
		agg_peserta_ujian AS (
			SELECT
				tj.id_jadwal_ujian,
				COUNT(pu.id_peserta_ujian)::integer AS total_peserta_ujian
			FROM target_jadwal tj
			JOIN peserta_ujian pu
				ON pu.id_jadwal_ujian = tj.id_jadwal_ujian
			GROUP BY tj.id_jadwal_ujian
		)
		INSERT INTO statistik_ujian (
			id_jadwal_ujian,
			nilai_tertinggi,
			nilai_terendah,
			nilai_rata_rata,
			total_peserta_ujian
		)
		SELECT
			ahu.id_jadwal_ujian,
			ahu.nilai_tertinggi,
			ahu.nilai_terendah,
			ahu.nilai_rata_rata,
			apu.total_peserta_ujian
		FROM agg_hasil_ujian ahu
		JOIN agg_peserta_ujian apu
			ON apu.id_jadwal_ujian = ahu.id_jadwal_ujian
		ON CONFLICT (id_jadwal_ujian) DO UPDATE SET
			nilai_tertinggi = EXCLUDED.nilai_tertinggi,
			nilai_terendah = EXCLUDED.nilai_terendah,
			nilai_rata_rata = EXCLUDED.nilai_rata_rata,
			total_peserta_ujian = EXCLUDED.total_peserta_ujian,
			updated_at = NOW()
	`

	tag, err := r.q.Exec(ctx, query, idAttempt)
	if err != nil {
		if mappedErr := mapGradingConstraintError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed upsert statistik ujian", "layer", "repo.db", "op", "ujian.grading.upsert_statistik_ujian.exec", "attempt_id", idAttempt, "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		r.loggerFor(ctx).Error(ctx, "failed upsert statistik ujian", "layer", "repo.db", "op", "ujian.grading.upsert_statistik_ujian.no_rows", "attempt_id", idAttempt, "err", coreerror.ErrNotFound)
		return coreerror.ErrNotFound
	}

	return nil
}
