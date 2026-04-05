package attemptrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type attemptScanner interface {
	Scan(dest ...any) error
}

func scanAttemptUjianRow(row attemptScanner) (ujian.AttemptUjian, error) {
	var (
		item        ujian.AttemptUjian
		status      string
		waktuMulai  sql.NullTime
		waktuSubmit sql.NullTime
		deadlineAt  sql.NullTime
	)

	err := row.Scan(
		&item.IdAttempt,
		&item.IdPesertaUjian,
		&status,
		&waktuMulai,
		&waktuSubmit,
		&deadlineAt,
	)
	if err != nil {
		return ujian.AttemptUjian{}, err
	}

	item.StatusAttempt = ujian.StatusAttempt(status)
	item.WaktuMulai = nullTimeToPtr(waktuMulai)
	item.WaktuSubmit = nullTimeToPtr(waktuSubmit)
	item.DeadlineAt = nullTimeToPtr(deadlineAt)

	return item, nil
}

func (r *AttemptUjianRepo) scanPesertaUjianSubmittedRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.PesertaUjianSubmitted, error) {
	results := make([]ujian.PesertaUjianSubmitted, 0)

	for rows.Next() {
		var (
			item        ujian.PesertaUjianSubmitted
			nilaiAkhir  sql.NullFloat64
			waktuMulai  sql.NullTime
			waktuSubmit sql.NullTime
		)

		if err := rows.Scan(
			&item.IdPesertaUjian,
			&item.IdAttempt,
			&item.IdSiswa,
			&item.TingkatKelas,
			&item.NamaKelas,
			&item.NamaLengkap,
			&item.NoAbsen,
			&nilaiAkhir,
			&waktuMulai,
			&waktuSubmit,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning peserta ujian submitted", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}

		item.NilaiAkhir = nullFloat64ToPtr(nilaiAkhir)
		item.WaktuMulai = nullTimeToPtr(waktuMulai)
		item.WaktuSubmit = nullTimeToPtr(waktuSubmit)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating peserta ujian submitted", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}

func nullFloat64ToPtr(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}

	v := value.Float64
	return &v
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
