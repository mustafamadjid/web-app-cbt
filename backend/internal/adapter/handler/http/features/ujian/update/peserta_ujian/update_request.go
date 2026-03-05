package httpx

import "time"

type UpdatePesertaUjianRequest struct {
	IdJadwalUjian *int       `json:"id_jadwal_ujian"`
	IdSiswa       *int       `json:"id_siswa"`
	WaktuMulai    *time.Time `json:"waktu_mulai"`
	WaktuSubmit   *time.Time `json:"waktu_submit"`
	NilaiUjian    *float64   `json:"nilai_ujian"`
}
