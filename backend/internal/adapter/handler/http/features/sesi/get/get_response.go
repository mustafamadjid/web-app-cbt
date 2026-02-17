package httpx

type SesiResponse struct {
	IdSesi   int    `json:"id_sesi"`
	KodeSesi string `json:"kode_sesi"`
	NamaSesi string `json:"nama_sesi"`
}

type ListSesiResponse struct {
	Items []SesiResponse `json:"items"`
}
