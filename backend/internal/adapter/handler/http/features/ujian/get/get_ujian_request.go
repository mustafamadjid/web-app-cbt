package httpx

import "time"

type GetPesertaUjianRequest struct {
	IDPesertaUjian *int
	IDJadwalUjian  *int
	IDSiswa        *int
	WaktuMulai     *time.Time
	WaktuSubmit    *time.Time
	NilaiUjian     *float64
}

type GetJawabanUjianRequest struct {
	IDJawaban      *int
	IDPesertaUjian *int
	IDSoal         *int
	IDPilihan      *int
	JawabanEssay   *string
	IsBenar        *bool
	WaktuJawab     *time.Time
}
