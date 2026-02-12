package httpx

type UpdateMapelRequest struct {
	IdKelas   *int    `json:"id_kelas"`
	KodeMapel *string `json:"kode_mapel"`
	NamaMapel *string `json:"nama_mapel"`
	Deskripsi *string `json:"deskripsi"`
}
