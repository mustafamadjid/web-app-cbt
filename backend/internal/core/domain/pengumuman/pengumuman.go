package pengumuman

type ID int


type Pengumuman struct {
	IdPengumuman             ID
	IdPengguna               ID
	JudulPengumuman          string
	IsiPengumuman            string
	TanggalRilisPengumuman   string
	TanggalSelesaiPengumuman string
	DokumenPengumuman        string
}