package httpx

import "time"

type AttemptUjianRequest struct {
	IdSiswa       int       `json:"id_siswa"`
	IdJadwalUjian int       `json:"id_jadwal_ujian"`
	TokenUjian    string    `json:"token_ujian"`
	WaktuMulai    time.Time `json:"waktu_mulai"`
}
