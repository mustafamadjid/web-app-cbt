package parser_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"testing"

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractParagraphContents_SuperscriptAndMath(t *testing.T) {
	t.Parallel()

	docx := buildDocx(`<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math">
<w:body>
	<w:p>
		<w:r><w:t>x</w:t></w:r>
		<w:r><w:rPr><w:vertAlign w:val="superscript"/></w:rPr><w:t>2</w:t></w:r>
	</w:p>
	<w:p>
		<m:oMath>
			<m:f>
				<m:num><m:r><m:t>a</m:t></m:r></m:num>
				<m:den><m:r><m:t>b</m:t></m:r></m:den>
			</m:f>
		</m:oMath>
	</w:p>
</w:body></w:document>`)

	paragraphs, warnings, err := parser.ExtractParagraphContents(docx)
	require.NoError(t, err)
	require.Len(t, warnings, 0)
	require.Len(t, paragraphs, 2)

	first := paragraphs[0]
	require.Len(t, first.Blocks, 1)
	require.Len(t, first.Blocks[0].Children, 2)
	assert.Equal(t, "x", first.Blocks[0].Children[0].Text)
	assert.Equal(t, []content.Mark{content.MarkSup}, first.Blocks[0].Children[1].Marks)
	assert.Equal(t, "2", first.Blocks[0].Children[1].Text)

	second := paragraphs[1]
	require.Len(t, second.Blocks, 1)
	require.Len(t, second.Blocks[0].Children, 1)
	assert.Equal(t, "math", second.Blocks[0].Children[0].Type)
	assert.Equal(t, `\frac{a}{b}`, second.Blocks[0].Children[0].Latex)
}

func TestParseMarkersFromContent_Branches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		paragraphs  []content.RichContent
		docxData    []byte
		expectedErr string
		validate    func(t *testing.T, result []importsoal.ParsedSoal)
	}{
		{
			name:       "Branch 1 -> paragraf kosong dilewati dan default tanpa soal aktif diabaikan",
			paragraphs: richParagraphs("", "teks tanpa marker"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				assert.Empty(t, result)
			},
		},
		{
			name: "Branch 2 -> [Q:PG] memulai soal pilihan ganda dengan rich content",
			paragraphs: []content.RichContent{
				richParagraph(
					content.Inline{Type: "text", Text: "[Q:PG] Hitung x"},
					content.Inline{Type: "text", Text: "2", Marks: []content.Mark{content.MarkSup}},
				),
				richParagraph(
					content.Inline{Type: "text", Text: "[A] "},
					content.Inline{Type: "math", Latex: `\frac{a}{b}`, Display: "inline"},
				),
				richParagraph(content.Inline{Type: "text", Text: "[B] Jawaban biasa"}),
				richParagraph(content.Inline{Type: "text", Text: "[ANS] A"}),
				richParagraph(content.Inline{Type: "text", Text: "[W] 5"}),
			},
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, "pilihan_ganda", result[0].TipeSoal)
				assert.Equal(t, "Hitung x2", result[0].Pertanyaan)
				assert.Equal(t, `\frac{a}{b}`, result[0].Opsi[0].Isi)
				assert.Equal(t, `\frac{a}{b}`, result[0].Opsi[0].IsiContent.Blocks[0].Children[1].Latex)
				assert.True(t, result[0].Opsi[0].IsBenar)
				assert.Equal(t, 5.0, result[0].BobotSoal)
			},
		},
		{
			name:       "Branch 3 -> [Q:ESSAY] memulai soal essay",
			paragraphs: richParagraphs("[Q:ESSAY] Jelaskan fotosintesis", "[ANS] Proses membuat makanan"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, "essay", result[0].TipeSoal)
				assert.Equal(t, "Jelaskan fotosintesis", result[0].Pertanyaan)
				assert.Equal(t, "Proses membuat makanan", result[0].KunciJawaban)
			},
		},
		{
			name:        "Branch 4 -> opsi tanpa soal aktif mengembalikan error",
			paragraphs:  richParagraphs("[A] Jawaban A"),
			expectedErr: "opsi ditemukan tanpa soal aktif",
		},
		{
			name: "Branch 5 -> marker opsi enumerasi Word tetap terbaca",
			paragraphs: richParagraphs(
				"[Q:PG] Pilih jawaban",
				"A. [A] Jawaban A",
				"B. [B] Jawaban B",
				"[ANS] A",
			),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				require.Len(t, result[0].Opsi, 2)
				assert.Equal(t, "A", result[0].Opsi[0].Label)
				assert.Equal(t, "Jawaban A", result[0].Opsi[0].Isi)
				assert.Equal(t, "B", result[0].Opsi[1].Label)
				assert.Equal(t, "Jawaban B", result[0].Opsi[1].Isi)
				require.NoError(t, parser.ValidateParsedSoal(result))
			},
		},
		{
			name:        "Branch 6 -> [ANS] tanpa soal aktif mengembalikan error",
			paragraphs:  richParagraphs("[ANS] A"),
			expectedErr: "[ANS] ditemukan tanpa soal aktif",
		},
		{
			name: "Branch 7 -> lanjutan jawaban digabung dengan baris sebelumnya",
			paragraphs: richParagraphs(
				"[Q:ESSAY] Jelaskan",
				"[ANS] Baris pertama",
				"Baris kedua",
			),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, "Baris pertama\nBaris kedua", result[0].KunciJawaban)
			},
		},
		{
			name:        "Branch 8 -> [IMG] tanpa soal aktif mengembalikan error",
			paragraphs:  richParagraphs("[IMG] image1.png"),
			expectedErr: "[IMG] ditemukan tanpa soal aktif",
		},
		{
			name:       "Branch 9 -> [IMG] eksplisit mengisi gambar soal",
			paragraphs: richParagraphs("[Q:PG] Perhatikan gambar", "[IMG] image1.png", "[A] A", "[B] B", "[ANS] A"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, "image1.png", result[0].Gambar)
			},
		},
		{
			name:       "Branch 10 -> [IMG] kosong memakai antrean gambar docx untuk soal",
			paragraphs: richParagraphs("[Q:PG] Perhatikan gambar", "[IMG]", "[A] A", "[B] B", "[ANS] A"),
			docxData:   buildDocxWithImages("image1.png"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, "image1.png", result[0].Gambar)
				assert.Empty(t, result[0].Opsi[0].IsiContent.Blocks[0].Children[0].Src)
			},
		},
		{
			name: "Branch 11 -> [IMG] setelah opsi masuk ke opsi terakhir",
			paragraphs: richParagraphs(
				"[Q:PG] Pilih gambar",
				"A. [A]",
				"[IMG]",
				"B. [B]",
				"[IMG]",
				"[ANS] A",
			),
			docxData: buildDocxWithImages("image1.png", "image2.png"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				require.Len(t, result[0].Opsi, 2)
				assert.Empty(t, result[0].Gambar)
				assert.Equal(t, "image1.png", result[0].Opsi[0].IsiContent.Blocks[0].Children[0].Src)
				assert.Equal(t, "image2.png", result[0].Opsi[1].IsiContent.Blocks[0].Children[0].Src)
				require.NoError(t, parser.ValidateParsedSoal(result))
			},
		},
		{
			name:        "Branch 12 -> [IMG] kosong tanpa sisa gambar mengembalikan error",
			paragraphs:  richParagraphs("[Q:PG] Perhatikan gambar", "[IMG]"),
			expectedErr: "[IMG] tetapi tidak ada gambar tersisa di dokumen",
		},
		{
			name:        "Branch 13 -> [W] tanpa soal aktif mengembalikan error",
			paragraphs:  richParagraphs("[W] 5"),
			expectedErr: "[W] ditemukan tanpa soal aktif",
		},
		{
			name:       "Branch 14 -> [W] inline mengisi bobot soal",
			paragraphs: richParagraphs("[Q:ESSAY] Jelaskan", "[W] 2.5"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, 2.5, result[0].BobotSoal)
			},
		},
		{
			name:       "Branch 15 -> lanjutan [W] mengisi bobot soal",
			paragraphs: richParagraphs("[Q:ESSAY] Jelaskan", "[W]", "3"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, 3.0, result[0].BobotSoal)
			},
		},
		{
			name:        "Branch 16 -> bobot tidak valid mengembalikan error",
			paragraphs:  richParagraphs("[Q:PG] Soal", "[W] abc"),
			expectedErr: "bobot bukan angka",
		},
		{
			name: "Branch 17 -> lanjutan default menambah teks pertanyaan",
			paragraphs: richParagraphs(
				"[Q:ESSAY] Jelaskan",
				"dengan lengkap",
			),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, "Jelaskan\ndengan lengkap", result[0].Pertanyaan)
			},
		},
		{
			name: "Branch 18 -> lanjutan default menambah teks opsi terakhir",
			paragraphs: richParagraphs(
				"[Q:PG] Pilih",
				"[A] Baris pertama",
				"Baris kedua",
				"[B] Pembanding",
				"[ANS] A",
			),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				require.Len(t, result[0].Opsi, 2)
				assert.Equal(t, "Baris pertama\nBaris kedua", result[0].Opsi[0].Isi)
			},
		},
		{
			name: "Branch 19 -> image opsi dengan caption memakai auto image",
			paragraphs: richParagraphs(
				"[Q:PG] Pilih gambar yang benar",
				"[A] [IMG] Gambar pilihan",
				"[B] Jawaban biasa",
				"[ANS] A",
			),
			docxData: buildDocxWithImages("image1.png"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				option := result[0].Opsi[0]
				assert.Equal(t, "Gambar pilihan", option.Isi)
				require.Len(t, option.IsiContent.Blocks, 2)
				assert.Equal(t, "image", option.IsiContent.Blocks[0].Children[0].Type)
				assert.Equal(t, "image1.png", option.IsiContent.Blocks[0].Children[0].Src)
				assert.Equal(t, "Gambar pilihan", option.IsiContent.Blocks[0].Children[0].Alt)
				assert.Equal(t, "Gambar pilihan", option.IsiContent.Blocks[1].Children[0].Text)
				assert.True(t, option.IsBenar)
			},
		},
		{
			name: "Branch 20 -> image opsi dengan filename eksplisit memakai gambar yang disebut",
			paragraphs: richParagraphs(
				"[Q:PG] Pilih gambar yang benar",
				"[A] [IMG] image2.jpg Gambar pilihan",
				"[B] Jawaban biasa",
				"[ANS] A",
			),
			docxData: buildDocxWithImages("image1.png", "image2.jpg"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				option := result[0].Opsi[0]
				assert.Equal(t, "Gambar pilihan", option.Isi)
				require.Len(t, option.IsiContent.Blocks, 2)
				assert.Equal(t, "image2.jpg", option.IsiContent.Blocks[0].Children[0].Src)
				assert.Equal(t, "Gambar pilihan", option.IsiContent.Blocks[1].Children[0].Text)
			},
		},
		{
			name: "Branch 21 -> image opsi eksplisit mengonsumsi antrean saat filename cocok item berikutnya",
			paragraphs: richParagraphs(
				"[Q:PG] Pilih gambar yang benar",
				"[A] [IMG] image1.png Gambar A",
				"[B] [IMG] Gambar B",
				"[ANS] A",
			),
			docxData: buildDocxWithImages("image1.png", "image2.png"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Equal(t, "image1.png", result[0].Opsi[0].IsiContent.Blocks[0].Children[0].Src)
				assert.Equal(t, "image2.png", result[0].Opsi[1].IsiContent.Blocks[0].Children[0].Src)
			},
		},
		{
			name: "Branch 22 -> image opsi tanpa caption tetap valid",
			paragraphs: richParagraphs(
				"[Q:PG] Pilih gambar yang benar",
				"[A] [IMG]",
				"[B] Jawaban biasa",
				"[ANS] A",
			),
			docxData: buildDocxWithImages("image1.png"),
			validate: func(t *testing.T, result []importsoal.ParsedSoal) {
				require.Len(t, result, 1)
				assert.Empty(t, result[0].Opsi[0].Isi)
				assert.Equal(t, "image1.png", result[0].Opsi[0].IsiContent.Blocks[0].Children[0].Src)
				require.NoError(t, parser.ValidateParsedSoal(result))
			},
		},
		{
			name: "Branch 23 -> image opsi filename eksplisit yang hilang mengembalikan error",
			paragraphs: richParagraphs(
				"[Q:PG] Pilih gambar yang benar",
				"[A] [IMG] image2.jpg Gambar pilihan",
				"[B] Jawaban biasa",
				"[ANS] A",
			),
			docxData:    buildDocxWithImages("image1.png"),
			expectedErr: `gambar "image2.jpg" tidak ditemukan`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			docxData := tc.docxData
			if docxData == nil {
				docxData = buildDocx(wrapXML())
			}

			result, warnings, err := parser.ParseMarkersFromContent(tc.paragraphs, docxData)
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Empty(t, warnings)
			if tc.validate != nil {
				tc.validate(t, result)
			}
		})
	}
}

func richParagraphs(items ...string) []content.RichContent {
	out := make([]content.RichContent, 0, len(items))
	for _, item := range items {
		out = append(out, richParagraph(content.Inline{Type: "text", Text: item}))
	}
	return out
}

func richParagraph(children ...content.Inline) content.RichContent {
	return content.RichContent{
		Version: 1,
		Blocks: []content.Block{
			{
				Type:     "paragraph",
				Children: children,
			},
		},
	}
}

func buildDocxWithImages(imageNames ...string) []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	var document bytes.Buffer
	document.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<w:body>`)
	for i := range imageNames {
		document.WriteString(`<w:p><w:r><w:drawing><a:blip r:embed="rId`)
		document.WriteString(fmt.Sprintf("%d", i+1))
		document.WriteString(`"/></w:drawing></w:r></w:p>`)
	}
	document.WriteString(`</w:body></w:document>`)

	fw, _ := zw.Create("word/document.xml")
	_, _ = fw.Write(document.Bytes())

	var rels bytes.Buffer
	rels.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for i, imageName := range imageNames {
		rels.WriteString(`<Relationship Id="rId`)
		rels.WriteString(fmt.Sprintf("%d", i+1))
		rels.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/`)
		rels.WriteString(imageName)
		rels.WriteString(`"/>`)
	}
	rels.WriteString(`</Relationships>`)
	relFile, _ := zw.Create("word/_rels/document.xml.rels")
	_, _ = relFile.Write(rels.Bytes())

	for _, imageName := range imageNames {
		media, _ := zw.Create("word/media/" + imageName)
		_, _ = media.Write([]byte("fake image"))
	}

	_ = zw.Close()
	return buf.Bytes()
}
