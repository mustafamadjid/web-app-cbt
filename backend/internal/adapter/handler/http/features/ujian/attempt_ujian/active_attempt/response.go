package httpx

type GetActiveAttemptUjianResponse struct {
	IDAttempt      int     `json:"id_attempt"`
	IDPesertaUjian int     `json:"id_peserta_ujian"`
	StatusAttempt  string  `json:"status_attempt"`
	WaktuMulai     *string `json:"waktu_mulai"`
	WaktuSubmit    *string `json:"waktu_submit"`
	DeadlineAt     *string `json:"deadline_at"`
}
