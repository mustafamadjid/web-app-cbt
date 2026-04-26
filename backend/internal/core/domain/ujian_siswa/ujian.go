package ujian

import "time"
import content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"

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
	StatusUjian   *StatusUjian

	IdAttempt ID

	IdPengawas       ID
	NamaPengawas     string
	PengawasUsername string

	PengawasNamaLengkap string

	IdSesi   ID
	NamaSesi string

	IdRuangan   ID
	NamaRuangan string

	DeskripsiUjian *string
	Token          string
	AcakSoal       bool
}

type SoalUjianSiswa struct {
	IdSoal            ID
	IdBankSoalVersion ID
	TipeSoal          string
	Pertanyaan        string
	PertanyaanContent content.RichContent
	Gambar            string
	BobotSoal         float64
	NoUrutSoal        int

	OpsiJawaban []OpsiPilganUjian
}

type OpsiPilganUjian struct {
	IdPilihanGanda    ID
	IdSoal            ID
	IsiPilihan        string
	IsiPilihanContent content.RichContent
	IsBenar           bool
}

func (status StatusUjian) ValidStatus() bool {
	switch status {
	case BELUM_MULAI, MULAI, SELESAI, DIBATALKAN:
		return true
	default:
		return false
	}
}
