package httpx

type GetJawabanUjianResponse struct {
	IDAttempt int                           `json:"id_attempt"`
	Jawaban   []GetJawabanUjianItemResponse `json:"jawaban"`
}

type GetJawabanUjianItemResponse struct {
	IDJawaban    int     `json:"id_jawaban"`
	IDSoal       int     `json:"id_soal"`
	IDPilihan    *int    `json:"id_pilihan"`
	JawabanEssay *string `json:"jawaban_essay"`
	WaktuJawab   *string `json:"waktu_jawab"`
}
