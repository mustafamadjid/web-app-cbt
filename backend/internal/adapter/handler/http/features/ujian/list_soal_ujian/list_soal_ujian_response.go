package httpx

type ListSoalUjianResponse struct {
	IDSoal            int                                `json:"id_soal"`
	IDBankSoalVersion int                                `json:"id_bank_soal_version"`
	TipeSoal          string                             `json:"tipe_soal"`
	Pertanyaan        string                             `json:"pertanyaan"`
	Gambar            string                             `json:"gambar"`
	BobotSoal         int                                `json:"bobot_soal"`
	NoUrutSoal        int                                `json:"no_urut_soal"`
	OpsiJawaban       []ListSoalUjianOpsiJawabanResponse `json:"opsi_jawaban"`
}

type ListSoalUjianOpsiJawabanResponse struct {
	IDPilihanGanda int    `json:"id_pilihan_ganda"`
	IDSoal         int    `json:"id_soal"`
	IsiPilihan     string `json:"isi_pilihan"`
	IsBenar        bool   `json:"is_benar"`
}
