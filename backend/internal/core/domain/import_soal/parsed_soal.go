package importsoal

type ParsedSoal struct {
	Pertanyaan   string
	TipeSoal     string // "pilihan_ganda" atau "essay"
	Gambar       string // path relatif gambar (opsional)
	BobotSoal    int
	NoUrutSoal   int
	Opsi         []ParsedOpsi // kosong untuk essay
	KunciJawaban string       // huruf opsi untuk PG, teks jawaban untuk essay
}

type ParsedOpsi struct {
	Label   string // "A", "B", "C", "D", "E"
	Isi     string
	IsBenar bool
}
