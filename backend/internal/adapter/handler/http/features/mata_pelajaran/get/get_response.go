package httpx

type MapelResponse struct {
	IDMapel   int    `json:"id_mapel"`
	IDKelas   int    `json:"id_kelas"`
	KodeMapel string `json:"kode_mapel"`
	NamaMapel string `json:"nama_mapel"`
	Deskripsi string `json:"deskripsi"`
}

type ListMapelResponse struct {
	Items []MapelResponse `json:"items"`
}
