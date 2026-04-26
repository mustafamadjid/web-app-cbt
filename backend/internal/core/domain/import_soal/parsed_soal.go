package importsoal

import content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"

type ParsedSoal struct {
	Pertanyaan        string
	PertanyaanContent content.RichContent
	TipeSoal          string // "pilihan_ganda" atau "essay"
	Gambar            string // path relatif gambar (opsional)
	BobotSoal         float64
	NoUrutSoal        int
	Opsi              []ParsedOpsi // kosong untuk essay
	KunciJawaban      string       // huruf opsi untuk PG, teks jawaban untuk essay
}

type ParsedOpsi struct {
	Label      string // "A", "B", "C", "D", "E"
	Isi        string
	IsiContent content.RichContent
	IsBenar    bool
}
