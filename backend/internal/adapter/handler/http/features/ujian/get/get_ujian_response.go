package httpx

type ListUjianByIdResponse struct {
	ID               int    `json:"id"`
	IDUjian          int    `json:"id_ujian"`
	IDGuru           int    `json:"id_guru"`
	IDPengawas       int    `json:"id_pengawas"`
	NamaUjian        string `json:"nama_ujian"`
	PengawasUjian    string `json:"pengawas_ujian"`
	TglUjian         string `json:"tgl_ujian"`
	TanggalUjian     string `json:"tanggal_ujian"`
	WaktuMulai       string `json:"waktu_mulai"`
	WaktuSelesai     string `json:"waktu_selesai"`
	SesiUjian        int    `json:"sesi_ujian"`
	RuangUjian       string `json:"ruang_ujian"`
	IDRuang          int    `json:"id_ruang"`
	StatusUjian      string `json:"status_ujian"`
	Started          int    `json:"started"`
	TingkatKelas     int    `json:"tingkat_kelas"`
	TingkatKelasID   int    `json:"tingkat_kelas_id"`
	NamaKelas        string `json:"nama_kelas"`
	PembuatUsername  string `json:"pembuat_username"`
	PengawasUsername string `json:"pengawas_username"`
	DeskripsiUjian   string `json:"deskripsi_ujian"`
	Token            string `json:"token"`
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

