package httpx

type CreateTingkatKelasReq struct {
	TingkatKelas int `json:"tingkat_kelas"`
}

type CreateNamaKelasReq struct {
	IdTingkatKelas int    `json:"id_tingkat_kelas"`
	NamaKelas      string `json:"nama_kelas"`
}
