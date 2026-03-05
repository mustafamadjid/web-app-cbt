package httpx

import "time"

type CreatePesertaUjianRequest struct {
	IdJadwalUjian int        `json:"id_jadwal_ujian"`
	IdSiswa       int        `json:"id_siswa"`
	WaktuMulai    *time.Time `json:"waktu_mulai,omitempty"`
}