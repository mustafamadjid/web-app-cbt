package ujian

import "time"

type ID int
type StatusUjian string

const (
	BELUM_MULAI StatusUjian = "BELUM_MULAI"
	MULAI       StatusUjian = "MULAI"
	SELESAI     StatusUjian = "SELESAI"
	DIBATALKAN  StatusUjian = "DIBATALKAN"
)

type Ujian struct {
	IdUjian        ID
	IdBankSoal     ID
	IdKelas        ID
	IdNamaKelas    *ID
	IdGuru         ID
	NamaUjian      string
	DeskripsiUjian *string
	AcakSoal       bool
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type JadwalUjian struct {
	IdJadwalUjian ID
	IdUjian       ID
	IdSesi        ID
	IdRuangan     ID
	TanggalUjian  time.Time
	WaktuMulai    time.Time
	WaktuSelesai  time.Time
	StatusUjian   StatusUjian
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type PesertaUjian struct {
	IdPesertaUjian ID
	IdJadwalUjian  ID
	IdSiswa        ID
	WaktuMulai     *time.Time
	WaktuSubmit    *time.Time
	NilaiUjian     *float64
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type JawabanUjianSiswa struct {
	IdJawaban      ID
	IdPesertaUjian ID
	IdSoal         ID
	IdPilihan      *ID
	JawabanEssay   *string
	IsBenar        *bool
	WaktuJawab     *time.Time
}

func (status StatusUjian) ValidStatus() bool {
	switch status {
	case BELUM_MULAI, MULAI, SELESAI, DIBATALKAN:
		return true
	default:
		return false
	}
}
