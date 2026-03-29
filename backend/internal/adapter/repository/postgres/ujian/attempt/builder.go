package attemptrepo

import (
	"database/sql"
	"time"

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
