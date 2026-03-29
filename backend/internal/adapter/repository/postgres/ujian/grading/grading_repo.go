package gradingrepo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
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

func (r *GradingRepo) InsertNilaiToHasilUjian(ctx context.Context, totalNilai float64, idAttempt ujian.ID) error {
	const updateQuery = `
		UPDATE hasil_ujian
		SET
			nilai_akhir = $1,
			graded_at = NOW()
		WHERE id_attempt = $2
	`

	tag, err := r.q.Exec(ctx, updateQuery, totalNilai, idAttempt)
	if err != nil {
		if mappedErr := mapGradingConstraintError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed update nilai hasil ujian before insert", "layer", "repo.db", "op", "ujian.grading.insert_nilai.update_existing", "attempt_id", idAttempt, "err", err)
		return err
	}

	if tag.RowsAffected() > 0 {
		return nil
	}

	const insertQuery = `
		INSERT INTO hasil_ujian (
			id_attempt,
			nilai_akhir,
			graded_at
		)
		VALUES ($1, $2, NOW())
	`

	_, err = r.q.Exec(ctx, insertQuery, idAttempt, totalNilai)
	if err != nil {
		if mappedErr := mapGradingConstraintError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed insert nilai hasil ujian", "layer", "repo.db", "op", "ujian.grading.insert_nilai.insert", "attempt_id", idAttempt, "err", err)
		return err
	}

	return nil
}

func (r *GradingRepo) UpdateNilaiInHasilUjian(ctx context.Context, totalNilai float64, idAttempt ujian.ID) error {
	const query = `
		UPDATE hasil_ujian
		SET
			nilai_akhir = $1,
			graded_at = NOW()
		WHERE id_attempt = $2
	`

	tag, err := r.q.Exec(ctx, query, totalNilai, idAttempt)
	if err != nil {
		if mappedErr := mapGradingConstraintError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed update nilai hasil ujian", "layer", "repo.db", "op", "ujian.grading.update_nilai", "attempt_id", idAttempt, "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
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

func buildStatistikSoalPayload(items []ujian.StatistikSoal, isBenar bool) []statistikSoalUpsertItem {
	if len(items) == 0 {
		return nil
	}

	payload := make([]statistikSoalUpsertItem, 0, len(items))

	for _, item := range items {
		entry := statistikSoalUpsertItem{
			IDSoal:  item.IDSoal,
			IDUjian: item.IDUjian,
		}

		if isBenar {
			entry.JumlahJawabanBenar++
		} else {
			entry.JumlahJawabanSalah++
		}

		payload = append(payload, entry)
	}

	return payload
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
