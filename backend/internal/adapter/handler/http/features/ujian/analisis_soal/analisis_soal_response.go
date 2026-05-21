package httpx

import content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"

type GetListAnalisisSoalResponse struct {
	IDJadwalUjian int                        `json:"id_jadwal_ujian"`
	AnalisisSoal  []AnalisisSoalItemResponse `json:"analisis_soal"`
}

type AnalisisSoalItemResponse struct {
	IDSoal             int                        `json:"id_soal"`
	IDBankSoalVersion  int                        `json:"id_bank_soal_version"`
	TipeSoal           string                     `json:"tipe_soal"`
	Pertanyaan         string                     `json:"pertanyaan"`
	PertanyaanContent  *content.RichContent       `json:"pertanyaan_content,omitempty"`
	Gambar             string                     `json:"gambar"`
	BobotSoal          float64                    `json:"bobot_soal"`
	NoUrutSoal         int                        `json:"no_urut_soal"`
	JumlahJawabanBenar int                        `json:"jumlah_jawaban_benar"`
	JumlahJawabanSalah int                        `json:"jumlah_jawaban_salah"`
	OpsiJawaban        []AnalisisSoalOpsiResponse `json:"opsi_jawaban"`
}

type AnalisisSoalOpsiResponse struct {
	IDPilihanGanda    int                  `json:"id_pilihan_ganda"`
	IsiPilihan        string               `json:"isi_pilihan"`
	IsiPilihanContent *content.RichContent `json:"isi_pilihan_content,omitempty"`
	IsBenar           bool                 `json:"is_benar"`
}
