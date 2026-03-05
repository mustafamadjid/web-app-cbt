package updatepatch

import (
	"time"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type UpdateUjianPatch struct {
	IdBankSoal     *ujian.ID
	IdKelas        *ujian.ID
	IdNamaKelas    *ujian.ID
	IdGuru         *ujian.ID
	NamaUjian      *string
	DeskripsiUjian *string
	AcakSoal       *bool
	UpdatedAt      *time.Time
}

type UpdateJadwalUjianPatch struct {
	IdUjian      *ujian.ID
	IdSesi       *ujian.ID
	IdRuangan    *ujian.ID
	IdPengawas   *ujian.ID
	TanggalUjian *time.Time
	Token        *string
	WaktuMulai   *time.Time
	WaktuSelesai *time.Time
	StatusUjian  *ujian.StatusUjian
	UpdatedAt    *time.Time
}

type UpdatePenjadwalanUjian struct {
	Ujian UpdateUjianPatch
	JadwalUjian UpdateJadwalUjianPatch
}

type UpdatePesertaUjianPatch struct {
	IdJadwalUjian *ujian.ID
	IdSiswa       *ujian.ID
	WaktuMulai    *time.Time
	WaktuSubmit   *time.Time
	NilaiUjian    *float64
	UpdatedAt     *time.Time
}

type UpdateJawabanUjianSiswaPatch struct {
	IdPesertaUjian *ujian.ID
	IdSoal         *ujian.ID
	IdPilihan      *ujian.ID
	JawabanEssay   *string
	IsBenar        *bool
	WaktuJawab     *time.Time
}

