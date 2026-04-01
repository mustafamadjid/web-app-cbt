package ujian

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type StatusAttempt string

const (
	ATTEMPT_IN_PROGRESS StatusAttempt = "in_progress"
	ATTEMPT_SUBMITTED   StatusAttempt = "submitted"
	ATTEMPT_EXPIRED     StatusAttempt = "expired"
	ATTEMPT_CANCELLED   StatusAttempt = "cancelled"
)

type AttemptUjian struct {
	IdAttempt      ID
	IdPesertaUjian ID
	StatusAttempt  StatusAttempt
	WaktuMulai     *time.Time
	WaktuSubmit    *time.Time
	DeadlineAt     *time.Time
}

type PesertaUjian struct {
	IdPesertaUjian ID
	IdJadwalUjian  ID
	IdSiswa        ID
	DataSiswa      user.DataSiswa
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}


type JawabanUjian struct {
	IdJawaban    ID `json:"id_jawaban"`
	IdSoal       ID	`json:"id_soal"`
	IdPilihan    *ID `json:"id_pilihan"`
	JawabanEssay *string `json:"jawaban_essay"`
	WaktuJawab   *time.Time `json:"waktu_jawab"`
	IdAttempt   ID `json:"id_attempt"`
	EssayIsBenar *bool	`json:"essay_is_benar"`
}

func (status StatusAttempt) ValidStatus() bool {
	switch status {
	case ATTEMPT_IN_PROGRESS, ATTEMPT_SUBMITTED, ATTEMPT_EXPIRED, ATTEMPT_CANCELLED:
		return true
	default:
		return false
	}
}
