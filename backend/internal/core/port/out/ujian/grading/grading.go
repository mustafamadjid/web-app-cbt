package grading_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type GradingUjianRepository interface {
	InsertNilaiToHasilUjian(ctx context.Context, totalNilai float64, idAttempt ujian.ID) error
	UpdateNilaiInHasilUjian(ctx context.Context, totalNilai float64, idAttempt ujian.ID) error

	UpsertJawabanBenarToStatistikSoal(ctx context.Context, soalBenar []ujian.StatistikSoal) error
	UpsertJawabanSalahToStatistikSoal(ctx context.Context, soalSalah []ujian.StatistikSoal) error
}