package httpx

type PengumumanGetResponse struct {
	IdPengumuman             int    `json:"id_pengumuman"`
	IdPengguna               int    `json:"id_pengguna"`
	JudulPengumuman          string `json:"judul_pengumuman"`
	IsiPengumuman            string `json:"isi_pengumuman"`
	TanggalRilisPengumuman   string `json:"tanggal_rilis_pengumuman"`
	TanggalSelesaiPengumuman string `json:"tanggal_selesai_pengumuman"`
	DokumenPengumuman        string `json:"dokumen_pengumuman,omitempty"`
}
