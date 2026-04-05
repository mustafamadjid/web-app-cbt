package grading_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type GradingUjianRepository interface {
	UpsertNilaiToHasilUjian(ctx context.Context, totalNilai float64, hasilUjian ujian.HasilUjian) error

	UpsertJawabanBenarToStatistikSoal(ctx context.Context, soalBenar []ujian.StatistikSoal) error
	UpsertJawabanSalahToStatistikSoal(ctx context.Context, soalSalah []ujian.StatistikSoal) error

	UpdateAndGradingEssayUjian(ctx context.Context,jawabanSiswa []ujian.JawabanUjian, gradedBy ujian.ID) error

	UpsertToStatistikUjian(ctx context.Context, idAttempt ujian.ID) error
}
type ListGradingUjianRepository interface {
	ListUjianEssayUngraded(ctx context.Context, filter query.ListUjianEssayUngradedFilter)([]ujian.ListUjian,error)

	// ListStatistikUjianByIdJadwalUjian(ctx context.Context,idJadwalUjian ujian.ID)([]ujian.StatistikUjian,error)
}
