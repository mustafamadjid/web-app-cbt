package httpx

type GetStatistikUjianResponse struct {
	IDStatistikUjian int     `json:"id_statistik_ujian"`
	IDJadwalUjian    int     `json:"id_jadwal_ujian"`
	NilaiTertinggi   float64 `json:"nilai_tertinggi"`
	NilaiTerendah    float64 `json:"nilai_terendah"`
	RataRata         float64 `json:"rata_rata"`
	JumlahPeserta    int     `json:"jumlah_peserta"`
}
