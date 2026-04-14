package parser_test

import (
	"archive/zip"
	"bytes"
	"testing"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/parser"
	"github.com/stretchr/testify/assert"
)

// --- helper: build a minimal .docx (ZIP with word/document.xml) ---

func buildDocx(xmlContent string) []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	fw, _ := zw.Create("word/document.xml")
	fw.Write([]byte(xmlContent))

	rels, _ := zw.Create("word/_rels/document.xml.rels")
	rels.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`))

	zw.Close()
	return buf.Bytes()
}

func wrapXML(paragraphs ...string) string {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>`
	for _, p := range paragraphs {
		xml += `<w:p><w:r><w:t>` + p + `</w:t></w:r></w:p>`
	}
	xml += `</w:body></w:document>`
	return xml
}

// ============================================================
// Test ExtractParagraphs
// ============================================================

func TestExtractParagraphs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      []byte
		expectLen int
		expectErr bool
	}{
		{
			name:      "path 1 -> data bukan ZIP",
			data:      []byte("bukan zip file"),
			expectErr: true,
		},
		{
			name:      "path 2 -> ZIP tanpa document.xml",
			data:      buildDocxNoDocXML(),
			expectErr: true,
		},
		{
			name:      "path 3 -> happy path dengan 2 paragraf",
			data:      buildDocx(wrapXML("Hello", "World")),
			expectLen: 2,
			expectErr: false,
		},
		{
			name:      "path 4 -> paragraf kosong dilewati",
			data:      buildDocx(wrapXML("Hello", "", "World")),
			expectLen: 2,
			expectErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := parser.ExtractParagraphs(tc.data)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tc.expectLen)
			}
		})
	}
}

func buildDocxNoDocXML() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fw, _ := zw.Create("word/other.xml")
	fw.Write([]byte("<root/>"))
	zw.Close()
	return buf.Bytes()
}

// ============================================================
// Test ParseMarkers
// ============================================================

func TestParseMarkers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		paragraphs []string
		expectLen  int
		expectErr  bool
		validate   func(t *testing.T, soal []importsoal.ParsedSoal)
	}{
		{
			name:       "path 1 -> input kosong",
			paragraphs: []string{},
			expectLen:  0,
			expectErr:  false,
		},
		{
			name:       "path 2 -> opsi tanpa soal aktif",
			paragraphs: []string{"[A] Jawaban A"},
			expectErr:  true,
		},
		{
			name:       "path 3 -> [ANS] tanpa soal aktif",
			paragraphs: []string{"[ANS] A"},
			expectErr:  true,
		},
		{
			name:       "path 4 -> [IMG] tanpa soal aktif",
			paragraphs: []string{"[IMG] image.png"},
			expectErr:  true,
		},
		{
			name:       "path 5 -> [W] tanpa soal aktif",
			paragraphs: []string{"[W] 10"},
			expectErr:  true,
		},
		{
			name:       "path 6 -> bobot bukan angka",
			paragraphs: []string{"[Q:PG] Soal?", "[W] abc"},
			expectErr:  true,
		},
		{
			name:       "path 7 -> bobot nol atau negatif",
			paragraphs: []string{"[Q:PG] Soal?", "[W] 0"},
			expectErr:  true,
		},
		{
			name: "path 8 -> happy path PG lengkap",
			paragraphs: []string{
				"[Q:PG] Siapa presiden pertama?",
				"[A] Soekarno",
				"[B] Soeharto",
				"[C] Habibie",
				"[ANS] A",
				"[W] 10",
			},
			expectLen: 1,
			validate: func(t *testing.T, soal []importsoal.ParsedSoal) {
				s := soal[0]
				assert.Equal(t, "pilihan_ganda", s.TipeSoal)
				assert.Equal(t, "Siapa presiden pertama?", s.Pertanyaan)
				assert.Len(t, s.Opsi, 3)
				assert.Equal(t, "A", s.KunciJawaban)
				assert.True(t, s.Opsi[0].IsBenar)
				assert.False(t, s.Opsi[1].IsBenar)
				assert.Equal(t, 10.0, s.BobotSoal)
			},
		},
		{
			name: "path 9 -> happy path ESSAY",
			paragraphs: []string{
				"[Q:ESSAY] Jelaskan fotosintesis!",
				"[ANS] Proses konversi cahaya menjadi energi",
				"[W] 20",
			},
			expectLen: 1,
			validate: func(t *testing.T, soal []importsoal.ParsedSoal) {
				s := soal[0]
				assert.Equal(t, "essay", s.TipeSoal)
				assert.Equal(t, "Jelaskan fotosintesis!", s.Pertanyaan)
				assert.Len(t, s.Opsi, 0)
				assert.Equal(t, "Proses konversi cahaya menjadi energi", s.KunciJawaban)
				assert.Equal(t, 20.0, s.BobotSoal)
			},
		},
		{
			name: "path 10 -> multiple soal (PG + ESSAY)",
			paragraphs: []string{
				"[Q:PG] Soal 1?",
				"[A] A", "[B] B",
				"[ANS] B",
				"[W] 5",
				"[Q:ESSAY] Soal 2?",
				"[W] 15",
			},
			expectLen: 2,
			validate: func(t *testing.T, soal []importsoal.ParsedSoal) {
				assert.Equal(t, "pilihan_ganda", soal[0].TipeSoal)
				assert.Equal(t, "essay", soal[1].TipeSoal)
				assert.Equal(t, 5.0, soal[0].BobotSoal)
				assert.Equal(t, 15.0, soal[1].BobotSoal)
			},
		},
		{
			name: "path 11 -> soal dengan gambar",
			paragraphs: []string{
				"[Q:PG] Perhatikan gambar!",
				"[IMG] image1.png",
				"[A] X", "[B] Y",
				"[ANS] A",
				"[W] 5",
			},
			expectLen: 1,
			validate: func(t *testing.T, soal []importsoal.ParsedSoal) {
				assert.Equal(t, "image1.png", soal[0].Gambar)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := parser.ParseMarkers(tc.paragraphs, buildDocx(wrapXML()))
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tc.expectLen)
				if tc.validate != nil {
					tc.validate(t, result)
				}
			}
		})
	}
}

// ============================================================
// Test ValidateParsedSoal
// ============================================================

func TestValidateParsedSoal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		soalList  []importsoal.ParsedSoal
		expectErr bool
	}{
		{
			name:      "path 1 -> daftar kosong",
			soalList:  []importsoal.ParsedSoal{},
			expectErr: true,
		},
		{
			name: "path 2 -> pertanyaan kosong",
			soalList: []importsoal.ParsedSoal{
				{TipeSoal: "pilihan_ganda", Pertanyaan: "", BobotSoal: 1, Opsi: []importsoal.ParsedOpsi{{Label: "A", Isi: "x"}, {Label: "B", Isi: "y"}}, KunciJawaban: "A"},
			},
			expectErr: true,
		},
		{
			name: "path 3 -> bobot nol",
			soalList: []importsoal.ParsedSoal{
				{TipeSoal: "pilihan_ganda", Pertanyaan: "Q?", BobotSoal: 0, Opsi: []importsoal.ParsedOpsi{{Label: "A", Isi: "x"}, {Label: "B", Isi: "y"}}, KunciJawaban: "A"},
			},
			expectErr: true,
		},
		{
			name: "path 4 -> PG kurang dari 2 opsi",
			soalList: []importsoal.ParsedSoal{
				{TipeSoal: "pilihan_ganda", Pertanyaan: "Q?", BobotSoal: 1, Opsi: []importsoal.ParsedOpsi{{Label: "A", Isi: "x"}}, KunciJawaban: "A"},
			},
			expectErr: true,
		},
		{
			name: "path 5 -> PG kunci jawaban kosong",
			soalList: []importsoal.ParsedSoal{
				{TipeSoal: "pilihan_ganda", Pertanyaan: "Q?", BobotSoal: 1, Opsi: []importsoal.ParsedOpsi{{Label: "A", Isi: "x"}, {Label: "B", Isi: "y"}}, KunciJawaban: ""},
			},
			expectErr: true,
		},
		{
			name: "path 6 -> PG kunci jawaban tidak cocok",
			soalList: []importsoal.ParsedSoal{
				{TipeSoal: "pilihan_ganda", Pertanyaan: "Q?", BobotSoal: 1, Opsi: []importsoal.ParsedOpsi{{Label: "A", Isi: "x"}, {Label: "B", Isi: "y"}}, KunciJawaban: "C"},
			},
			expectErr: true,
		},
		{
			name: "path 7 -> tipe soal tidak valid",
			soalList: []importsoal.ParsedSoal{
				{TipeSoal: "unknown", Pertanyaan: "Q?", BobotSoal: 1},
			},
			expectErr: true,
		},
		{
			name: "path 8 -> happy path PG valid",
			soalList: []importsoal.ParsedSoal{
				{TipeSoal: "pilihan_ganda", Pertanyaan: "Q?", BobotSoal: 5, Opsi: []importsoal.ParsedOpsi{{Label: "A", Isi: "x"}, {Label: "B", Isi: "y"}}, KunciJawaban: "A"},
			},
			expectErr: false,
		},
		{
			name: "path 9 -> happy path essay valid",
			soalList: []importsoal.ParsedSoal{
				{TipeSoal: "essay", Pertanyaan: "Jelaskan!", BobotSoal: 10},
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := parser.ValidateParsedSoal(tc.soalList)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
