package httpx

import "time"

type UpdatePenjadwalanUjianRequest struct {
	// Ujian
	IdBankSoal     *int    `json:"id_bank_soal"`
	IdKelas        *int    `json:"id_kelas"`
	IdNamaKelas    *int    `json:"id_nama_kelas"`
	IdGuru         *int    `json:"id_guru"`
	NamaUjian      *string `json:"nama_ujian"`
	DeskripsiUjian *string `json:"deskripsi_ujian"`
	AcakSoal       *bool   `json:"acak_soal"`

	// Jadwal Ujian
	IdSesi       *int       `json:"id_sesi"`
	IdRuangan    *int       `json:"id_ruangan"`
	TanggalUjian *time.Time `json:"tanggal_ujian"`
	WaktuMulai   *time.Time `json:"waktu_mulai"`
	WaktuSelesai *time.Time `json:"waktu_selesai"`
	StatusUjian  *string    `json:"status_ujian"`
	Token        *string    `json:"token"`
	IdPengawas   *int       `json:"id_pengawas"`
}
