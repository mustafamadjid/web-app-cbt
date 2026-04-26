package httpx

import content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"

type HasilJawabanUjianResponse struct {
	IDAttempt    int                             `json:"id_attempt"`
	NilaiAkhir   *float64                        `json:"nilai_akhir"`
	HasilJawaban []HasilJawabanUjianItemResponse `json:"hasil_jawaban"`
}

type HasilJawabanUjianItemResponse struct {
	IDSoal            int                                    `json:"id_soal"`
	IDBankSoalVersion int                                    `json:"id_bank_soal_version"`
	TipeSoal          string                                 `json:"tipe_soal"`
	Pertanyaan        string                                 `json:"pertanyaan"`
	PertanyaanContent *content.RichContent                   `json:"pertanyaan_content,omitempty"`
	Gambar            string                                 `json:"gambar"`
	BobotSoal         float64                                `json:"bobot_soal"`
	NoUrutSoal        int                                    `json:"no_urut_soal"`
	OpsiJawaban       []HasilJawabanUjianOpsiJawabanResponse `json:"opsi_jawaban"`
	JawabanSiswa      HasilJawabanUjianJawabanSiswaResponse  `json:"jawaban_siswa"`
}

type HasilJawabanUjianOpsiJawabanResponse struct {
	IDPilihanGanda int `json:"id_pilihan_ganda"`
	// IDSoal         int    `json:"id_soal"`
	IsiPilihan        string               `json:"isi_pilihan"`
	IsiPilihanContent *content.RichContent `json:"isi_pilihan_content,omitempty"`
	IsBenar           bool                 `json:"is_benar"`
}

type HasilJawabanUjianJawabanSiswaResponse struct {
	IDJawaban    *int    `json:"id_jawaban"`
	IDPilihan    *int    `json:"id_pilihan"`
	JawabanEssay *string `json:"jawaban_essay"`
	WaktuJawab   *string `json:"waktu_jawab"`
	EssayIsBenar *bool   `json:"essay_is_benar"`
}
