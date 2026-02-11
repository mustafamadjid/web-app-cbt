package httpx

type UpdateKelasRequest struct {
	IdTingkatKelas *int    `json:"id_tingkat_kelas"`
	NamaKelas      *string `json:"nama_kelas"`
}