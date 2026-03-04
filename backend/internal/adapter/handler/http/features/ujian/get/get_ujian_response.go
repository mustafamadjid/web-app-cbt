package httpx

type UjianResponse struct {
	IDUjian        int     `json:"id_ujian"`
	IDBankSoal     int     `json:"id_bank_soal"`
	IDKelas        int     `json:"id_kelas"`
	IDNamaKelas    *int    `json:"id_nama_kelas"`
	IDGuru         int     `json:"id_guru"`
	NamaUjian      string  `json:"nama_ujian"`
	DeskripsiUjian *string `json:"deskripsi_ujian"`
	AcakSoal       bool    `json:"acak_soal"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      *string `json:"updated_at"`
}

type JadwalUjianResponse struct {
	IDJadwalUjian int     `json:"id_jadwal_ujian"`
	IDUjian       int     `json:"id_ujian"`
	IDSesi        int     `json:"id_sesi"`
	IDRuangan     int     `json:"id_ruangan"`
	IDPengawas    int     `json:"id_pengawas"`
	Token         string  `json:"token"`
	TanggalUjian  string  `json:"tanggal_ujian"`
	WaktuMulai    string  `json:"waktu_mulai"`
	WaktuSelesai  string  `json:"waktu_selesai"`
	StatusUjian   string  `json:"status_ujian"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
}

type PesertaUjianResponse struct {
	IDPesertaUjian int      `json:"id_peserta_ujian"`
	IDJadwalUjian  int      `json:"id_jadwal_ujian"`
	IDSiswa        int      `json:"id_siswa"`
	WaktuMulai     *string  `json:"waktu_mulai"`
	WaktuSubmit    *string  `json:"waktu_submit"`
	NilaiUjian     *float64 `json:"nilai_ujian"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      *string  `json:"updated_at"`
}

type JawabanUjianSiswaResponse struct {
	IDJawaban      int     `json:"id_jawaban"`
	IDPesertaUjian int     `json:"id_peserta_ujian"`
	IDSoal         int     `json:"id_soal"`
	IDPilihan      *int    `json:"id_pilihan"`
	JawabanEssay   *string `json:"jawaban_essay"`
	IsBenar        *bool   `json:"is_benar"`
	WaktuJawab     *string `json:"waktu_jawab"`
}
