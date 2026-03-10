package updatepatch

import (
	"time"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type UpdateAttemptUjianPatch struct {
	IdPesertaUjian *ujian.ID
	StatusAttempt  *ujian.StatusAttempt
	WaktuMulai     *time.Time
	WaktuSubmit    *time.Time
	DeadlineAt     *time.Time
	UpdatedAt      *time.Time
}

type UpdateHasilUjianPatch struct {
	IdAttempt  *ujian.ID
	GradedBy   *ujian.ID
	NilaiAkhir *float64
	Passed     *bool
	GradedAt   *time.Time
}

type UpdateJawabanUjianPatch struct {
	IdAttempt    *ujian.ID
	IdSoal       *ujian.ID
	IdPilihan    *ujian.ID
	JawabanEssay *string
	IsBenar      *bool
	WaktuJawab   *time.Time
}
