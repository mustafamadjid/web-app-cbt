package httpx

type ListMapelRequest struct {
	Search       string `json:"search"`
	TingkatKelas *int   `json:"tingkat_kelas"`
	NamaMapel    string `json:"nama_mapel"`
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
}
