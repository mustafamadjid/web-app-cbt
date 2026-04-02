package httpx

type ListUjianEssayUngradedResponse struct {
	ID               int    `json:"id"`
	IDUjian          int    `json:"id_ujian"`
	IDBankSoal       int    `json:"id_bank_soal"`
	IDGuru           int    `json:"id_guru"`
	IDPengawas       int    `json:"id_pengawas"`
	NamaUjian        string `json:"nama_ujian"`
	PengawasUjian    string `json:"pengawas_ujian"`
	TglUjian         string `json:"tgl_ujian"`
	TanggalUjian     string `json:"tanggal_ujian"`
	WaktuMulai       string `json:"waktu_mulai"`
	WaktuSelesai     string `json:"waktu_selesai"`
	IDSesi           int    `json:"id_sesi"`
	NamaSesi         string `json:"nama_sesi"`
	RuangUjian       string `json:"ruang_ujian"`
	IDRuang          int    `json:"id_ruang"`
	StatusUjian      string `json:"status_ujian"`
	Started          int    `json:"started"`
	TingkatKelas     int    `json:"tingkat_kelas"`
	TingkatKelasID   int    `json:"tingkat_kelas_id"`
	NamaKelas        string `json:"nama_kelas"`
	PembuatUsername  string `json:"pembuat_username"`
	PengawasUsername string `json:"pengawas_username"`
}
