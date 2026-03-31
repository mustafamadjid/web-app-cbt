package ujian

import "time"

type HasilUjian struct {
	IdHasil     ID
	IdAttempt   ID
	GradedBy    *ID
	NilaiAkhir  *float64
	Passed      *bool
	EssayGraded *bool
	GradedAt    *time.Time
}
