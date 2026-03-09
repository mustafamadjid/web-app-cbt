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
	Ujian       UpdateUjianPatch
	JadwalUjian UpdateJadwalUjianPatch
}
