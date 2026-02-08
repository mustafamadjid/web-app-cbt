package httpx

type TingkatKelasResponse struct {
	IDTingkatKelas int `json:"id_tingkat_kelas"`
	TingkatKelas   int `json:"tingkat_kelas"`
}

type NamaKelasResponse struct {
	IDNamaKelas    int    `json:"id_nama_kelas"`
	IDTingkatKelas int    `json:"id_tingkat_kelas"`
	NamaKelas      string `json:"nama_kelas"`
}

type FullKelasResponse struct {
	ItemsTingkatKelas []TingkatKelasResponse `json:"item_tingkat_kelas"`
	ItemsNamaKelas    []NamaKelasResponse    `json:"item_nama_kelas"`
}
