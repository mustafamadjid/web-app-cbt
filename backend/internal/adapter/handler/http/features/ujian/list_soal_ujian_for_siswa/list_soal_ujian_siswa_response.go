package httpx

import content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"

type ListSoalUjianSiswaResponse struct {
	IDSoal            int                                     `json:"id_soal"`
	TipeSoal          string                                  `json:"tipe_soal"`
	Pertanyaan        string                                  `json:"pertanyaan"`
	PertanyaanContent *content.RichContent                    `json:"pertanyaan_content,omitempty"`
	Gambar            string                                  `json:"gambar"`
	BobotSoal         float64                                 `json:"bobot_soal"`
	NoUrutSoal        int                                     `json:"no_urut_soal"`
	OpsiJawaban       []ListSoalUjianSiswaOpsiJawabanResponse `json:"opsi_jawaban"`
}

type ListSoalUjianSiswaOpsiJawabanResponse struct {
	IDPilihanGanda    int                  `json:"id_pilihan_ganda"`
	IsiPilihan        string               `json:"isi_pilihan"`
	IsiPilihanContent *content.RichContent `json:"isi_pilihan_content,omitempty"`
}
