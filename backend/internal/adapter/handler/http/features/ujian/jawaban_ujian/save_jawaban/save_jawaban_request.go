package httpx

import "time"

type SaveJawabanRequest struct {
	IDAttempt int                      `json:"id_attempt"`
	Jawaban   []SaveJawabanItemRequest `json:"jawaban"`
}

type SaveJawabanItemRequest struct {
	IDSoal       int        `json:"id_soal"`
	IDPilihan    *int       `json:"id_pilihan"`
	JawabanEssay *string    `json:"jawaban_essay"`
	WaktuJawab   *time.Time `json:"waktu_jawab"`
}
