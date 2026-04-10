package gradingrepo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func (r *GradingRepo) UpdateAndGradingEssayUjian(ctx context.Context, jawabanSiswa []ujian.JawabanUjian, gradedBy ujian.ID) error {
	if len(jawabanSiswa) == 0 {
		return nil
	}

	if r.pool == nil {
		err := errors.New("grading repo requires pgx pool for update transaction")
		r.loggerFor(ctx).Error(ctx, "failed update essay grading", "layer", "repo.db", "op", "ujian.grading.essay.precondition", "err", err)
		return err
	}

	payloadJSON, err := json.Marshal(jawabanSiswa)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed marshal essay grading payload", "layer", "repo.db", "op", "ujian.grading.essay.marshal_payload", "err", err)
		return err
	}

	// Transaction begin

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed begin tx essay grading", "layer", "repo.db", "op", "ujian.grading.essay.begin_tx", "err", err)
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	const updateQuery = `
		WITH data AS (
			SELECT *
			FROM jsonb_to_recordset($1::jsonb) AS x (
				id_jawaban bigint,
				essay_is_benar boolean
			)
		), target_rows AS (
			SELECT
				jus.id_jawaban,
				jus.id_attempt,
				jus.id_soal,
				d.essay_is_benar AS next_essay_is_benar
			FROM jawaban_ujian_siswa jus
			JOIN data d
				ON d.id_jawaban = jus.id_jawaban
			WHERE jus.essay_is_benar IS DISTINCT FROM d.essay_is_benar
			FOR UPDATE
		), updated AS (
			UPDATE jawaban_ujian_siswa jus
			SET essay_is_benar = t.next_essay_is_benar
			FROM target_rows t
			WHERE jus.id_jawaban = t.id_jawaban
			RETURNING
				jus.id_jawaban,
				t.id_attempt,
				t.id_soal,
				jus.essay_is_benar
		), affected_attempts AS (
			SELECT DISTINCT id_attempt
			FROM updated
		), affected_soal AS (
			SELECT DISTINCT
				u.id_soal,
				ju.id_ujian
			FROM updated u
			JOIN attempt_ujian au
				ON au.id_attempt = u.id_attempt
			JOIN peserta_ujian pu
				ON pu.id_peserta_ujian = au.id_peserta_ujian
			JOIN jadwal_ujian ju
				ON ju.id_jadwal_ujian = pu.id_jadwal_ujian
		), jawaban_current AS (
			SELECT
				jus.id_jawaban,
				jus.id_attempt,
				jus.id_soal,
				jus.id_pilihan,
				jus.jawaban_essay,
				CASE
					WHEN u.id_jawaban IS NOT NULL THEN u.essay_is_benar
					ELSE jus.essay_is_benar
				END AS essay_is_benar
			FROM jawaban_ujian_siswa jus
			LEFT JOIN updated u
				ON u.id_jawaban = jus.id_jawaban
			WHERE jus.id_attempt IN (SELECT id_attempt FROM affected_attempts)
				OR jus.id_soal IN (SELECT id_soal FROM affected_soal)
		), statistik_agg AS (
			SELECT
				a.id_soal,
				a.id_ujian,
				COUNT(*) FILTER (WHERE jc.essay_is_benar = true)::integer AS jumlah_jawaban_benar,
				COUNT(*) FILTER (WHERE jc.essay_is_benar = false)::integer AS jumlah_jawaban_salah
			FROM affected_soal a
			LEFT JOIN jawaban_current jc
				ON jc.id_soal = a.id_soal
			GROUP BY a.id_soal, a.id_ujian
		), statistik_upsert AS (
			INSERT INTO statistik_soal (
				id_soal,
				id_ujian,
				jumlah_jawaban_benar,
				jumlah_jawaban_salah
			)
			SELECT
				sa.id_soal,
				sa.id_ujian,
				sa.jumlah_jawaban_benar,
				sa.jumlah_jawaban_salah
			FROM statistik_agg sa
			ON CONFLICT (id_ujian, id_soal)
			DO UPDATE SET
				jumlah_jawaban_benar = EXCLUDED.jumlah_jawaban_benar,
				jumlah_jawaban_salah = EXCLUDED.jumlah_jawaban_salah,
				updated_at = NOW()
		), score_agg AS (
			SELECT
				aa.id_attempt,
				ju.id_jadwal_ujian,
				COALESCE(
					SUM(
						CASE
							WHEN LOWER(REPLACE(COALESCE(s.tipe_soal, ''), ' ', '_')) = 'essay'
								AND jc.essay_is_benar = true
							THEN s.bobot_soal
							WHEN LOWER(REPLACE(COALESCE(s.tipe_soal, ''), ' ', '_')) = 'pilihan_ganda'
								AND op.is_benar = true
							THEN s.bobot_soal
							ELSE 0
						END
					),
					0
				)::decimal(5,2) AS nilai_akhir,
				COUNT(*) FILTER (
					WHERE LOWER(REPLACE(COALESCE(s.tipe_soal, ''), ' ', '_')) = 'essay'
						AND COALESCE(NULLIF(BTRIM(jc.jawaban_essay), ''), '') <> ''
						AND jc.essay_is_benar IS NULL
				)::integer AS jumlah_essay_belum_dikoreksi
			FROM affected_attempts aa
			JOIN attempt_ujian au
				ON au.id_attempt = aa.id_attempt
			JOIN peserta_ujian pu
				ON pu.id_peserta_ujian = au.id_peserta_ujian
			JOIN jadwal_ujian ju
				ON ju.id_jadwal_ujian = pu.id_jadwal_ujian
			LEFT JOIN jawaban_current jc
				ON jc.id_attempt = aa.id_attempt
			LEFT JOIN isi_soal s
				ON s.id_soal = jc.id_soal
			LEFT JOIN opsi_pilihan_ganda op
				ON op.id_pilihan_ganda = jc.id_pilihan
			GROUP BY aa.id_attempt, ju.id_jadwal_ujian
		), hasil_upsert AS (
			INSERT INTO hasil_ujian (
				id_attempt,
				graded_by,
				nilai_akhir,
				essay_graded,
				graded_at,
				id_jadwal_ujian
			)
			SELECT
				sa.id_attempt,
				$2,
				sa.nilai_akhir,
				sa.jumlah_essay_belum_dikoreksi = 0,
				NOW(),
				sa.id_jadwal_ujian
			FROM score_agg sa
			ON CONFLICT (id_attempt)
			DO UPDATE SET
				graded_by = EXCLUDED.graded_by,
				nilai_akhir = EXCLUDED.nilai_akhir,
				essay_graded = EXCLUDED.essay_graded,
				graded_at = EXCLUDED.graded_at,
				id_jadwal_ujian = COALESCE(hasil_ujian.id_jadwal_ujian, EXCLUDED.id_jadwal_ujian)
			RETURNING id_jadwal_ujian
		), statistik_ujian_agg AS (
			SELECT
				hu.id_jadwal_ujian,
				MAX(hu.nilai_akhir)::decimal(5,2) AS nilai_tertinggi,
				MIN(hu.nilai_akhir)::decimal(5,2) AS nilai_terendah,
				ROUND(AVG(hu.nilai_akhir)::numeric, 2)::decimal(5,2) AS nilai_rata_rata,
				COUNT(*)::integer AS total_peserta_ujian
			FROM hasil_ujian hu
			WHERE hu.id_jadwal_ujian IN (SELECT id_jadwal_ujian FROM hasil_upsert)
				AND hu.nilai_akhir IS NOT NULL
			GROUP BY hu.id_jadwal_ujian
		), statistik_ujian_upsert AS (
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
				updated_at = NOW()
		)
		SELECT COUNT(*)::bigint
		FROM updated;
	`

	var updatedCount int64

	err = tx.QueryRow(ctx, updateQuery, payloadJSON, gradedBy).Scan(&updatedCount)
	if err != nil {
		if mappedErr := mapGradingConstraintError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed exec essay grading query", "layer", "repo.db", "op", "ujian.grading.essay.exec", "err", err)
		return err
	}

	if updatedCount == 0 {
		r.loggerFor(ctx).Error(ctx, "failed update essay grading", "layer", "repo.db", "op", "ujian.grading.essay.no_rows", "err", coreerror.ErrNotFound)
		return coreerror.ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit tx essay grading", "layer", "repo.db", "op", "ujian.grading.essay.commit_tx", "err", err)
		return err
	}

	committed = true
	return nil
}
