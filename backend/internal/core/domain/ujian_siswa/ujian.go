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
	Token         string
	IdPengawas    ID
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type PenjadwalanUjian struct {
	Ujian       Ujian
	JadwalUjian JadwalUjian
}

type ListUjian struct {
	IdUjian         ID
	IdBankSoal      ID
	IdGuru          ID
	NamaUjian       string
	PembuatUsername string

	IdKelas      ID
	IdNamaKelas  *ID
	TingkatKelas int
	NamaKelas    *string

	IdJadwalUjian ID
	TanggalUjian  time.Time
	WaktuMulai    time.Time
	WaktuSelesai  time.Time
	StatusUjian   StatusUjian

	IdPengawas       ID
	NamaPengawas     string
	PengawasUsername string

	IdSesi   ID
	NamaSesi string

	IdRuangan   ID
	NamaRuangan string

	DeskripsiUjian *string
	Token 			string
	AcakSoal   		bool
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

type SoalUjianSiswa struct {
	IdSoal ID
	IdBankSoalVersion ID
	TipeSoal string
	Pertanyaan string
	Gambar  string
	BobotSoal int
	NoUrutSoal int

	OpsiJawaban []OpsiPilganUjian
}

type OpsiPilganUjian struct {
	IdPilihanGanda ID
	IdSoal ID
	IsiPilihan string
	IsBenar bool
}

func (status StatusUjian) ValidStatus() bool {
	switch status {
	case BELUM_MULAI, MULAI, SELESAI, DIBATALKAN:
		return true
	default:
		return false
	}
}


