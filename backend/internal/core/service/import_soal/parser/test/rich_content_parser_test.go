package parser_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"testing"

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
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

func TestParseMarkersFromContent_RichQuestionAndOption(t *testing.T) {
	t.Parallel()

	paragraphs := []content.RichContent{
		{
			Version: 1,
			Blocks: []content.Block{
				{
					Type: "paragraph",
					Children: []content.Inline{
						{Type: "text", Text: "[Q:PG] Hitung x"},
						{Type: "text", Text: "2", Marks: []content.Mark{content.MarkSup}},
					},
				},
			},
		},
		{
			Version: 1,
			Blocks: []content.Block{
				{
					Type: "paragraph",
					Children: []content.Inline{
						{Type: "text", Text: "[A] "},
						{Type: "math", Latex: `\frac{a}{b}`, Display: "inline"},
					},
				},
			},
		},
		{
			Version: 1,
			Blocks: []content.Block{
				{
					Type:     "paragraph",
					Children: []content.Inline{{Type: "text", Text: "[B] Jawaban biasa"}},
				},
			},
		},
		{
			Version: 1,
			Blocks: []content.Block{
				{
					Type:     "paragraph",
					Children: []content.Inline{{Type: "text", Text: "[ANS] A"}},
				},
			},
		},
		{
			Version: 1,
			Blocks: []content.Block{
				{
					Type:     "paragraph",
					Children: []content.Inline{{Type: "text", Text: "[W] 5"}},
				},
			},
		},
	}

	result, warnings, err := parser.ParseMarkersFromContent(paragraphs, buildDocx(wrapXML()))
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, result, 1)

	assert.Equal(t, "Hitung x2", result[0].Pertanyaan)
	require.Len(t, result[0].PertanyaanContent.Blocks, 1)
	require.Len(t, result[0].Opsi, 2)
	assert.Equal(t, `\frac{a}{b}`, result[0].Opsi[0].Isi)
	assert.Equal(t, `\frac{a}{b}`, result[0].Opsi[0].IsiContent.Blocks[0].Children[1].Latex)
	assert.True(t, result[0].Opsi[0].IsBenar)
}

func TestParseMarkersFromContent_OptionImageAutoWithCaption(t *testing.T) {
	t.Parallel()

	paragraphs := richParagraphs(
		"[Q:PG] Pilih gambar yang benar",
		"[A] [IMG] Gambar pilihan",
		"[B] Jawaban biasa",
		"[ANS] A",
	)

	result, warnings, err := parser.ParseMarkersFromContent(paragraphs, buildDocxWithImages("image1.png"))
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, result, 1)
	require.Len(t, result[0].Opsi, 2)

	option := result[0].Opsi[0]
	assert.Equal(t, "Gambar pilihan", option.Isi)
	require.Len(t, option.IsiContent.Blocks, 2)
	require.Len(t, option.IsiContent.Blocks[0].Children, 1)
	assert.Equal(t, "image", option.IsiContent.Blocks[0].Children[0].Type)
	assert.Equal(t, "image1.png", option.IsiContent.Blocks[0].Children[0].Src)
	assert.Equal(t, "Gambar pilihan", option.IsiContent.Blocks[0].Children[0].Alt)
	assert.Equal(t, "Gambar pilihan", option.IsiContent.Blocks[1].Children[0].Text)
	assert.True(t, option.IsBenar)
}

func TestParseMarkersFromContent_OptionImageExplicitFilename(t *testing.T) {
	t.Parallel()

	paragraphs := richParagraphs(
		"[Q:PG] Pilih gambar yang benar",
		"[A] [IMG] image2.jpg Gambar pilihan",
		"[B] Jawaban biasa",
		"[ANS] A",
	)

	result, warnings, err := parser.ParseMarkersFromContent(paragraphs, buildDocxWithImages("image1.png", "image2.jpg"))
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, result, 1)
	require.Len(t, result[0].Opsi, 2)

	option := result[0].Opsi[0]
	assert.Equal(t, "Gambar pilihan", option.Isi)
	require.Len(t, option.IsiContent.Blocks, 2)
	assert.Equal(t, "image2.jpg", option.IsiContent.Blocks[0].Children[0].Src)
	assert.Equal(t, "Gambar pilihan", option.IsiContent.Blocks[1].Children[0].Text)
}

func TestParseMarkersFromContent_ExplicitOptionImageConsumesCurrentQueueItem(t *testing.T) {
	t.Parallel()

	paragraphs := richParagraphs(
		"[Q:PG] Pilih gambar yang benar",
		"[A] [IMG] image1.png Gambar A",
		"[B] [IMG] Gambar B",
		"[ANS] A",
	)

	result, warnings, err := parser.ParseMarkersFromContent(paragraphs, buildDocxWithImages("image1.png", "image2.png"))
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, result, 1)
	require.Len(t, result[0].Opsi, 2)

	assert.Equal(t, "image1.png", result[0].Opsi[0].IsiContent.Blocks[0].Children[0].Src)
	assert.Equal(t, "image2.png", result[0].Opsi[1].IsiContent.Blocks[0].Children[0].Src)
}

func TestParseMarkersFromContent_OptionImageOnlyIsValid(t *testing.T) {
	t.Parallel()

	paragraphs := richParagraphs(
		"[Q:PG] Pilih gambar yang benar",
		"[A] [IMG]",
		"[B] Jawaban biasa",
		"[ANS] A",
	)

	result, warnings, err := parser.ParseMarkersFromContent(paragraphs, buildDocxWithImages("image1.png"))
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, result, 1)
	require.Len(t, result[0].Opsi, 2)
	assert.Empty(t, result[0].Opsi[0].Isi)
	assert.Equal(t, "image1.png", result[0].Opsi[0].IsiContent.Blocks[0].Children[0].Src)
	require.NoError(t, parser.ValidateParsedSoal(result))
}

func TestParseMarkersFromContent_EnumeratedOptionMarkers(t *testing.T) {
	t.Parallel()

	paragraphs := richParagraphs(
		"[Q:PG] Pilih jawaban",
		"A. [A] Jawaban A",
		"B. [B] Jawaban B",
		"[ANS] A",
	)

	result, warnings, err := parser.ParseMarkersFromContent(paragraphs, buildDocxWithImages())
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, result, 1)
	require.Len(t, result[0].Opsi, 2)

	assert.Equal(t, "A", result[0].Opsi[0].Label)
	assert.Equal(t, "Jawaban A", result[0].Opsi[0].Isi)
	assert.Equal(t, "B", result[0].Opsi[1].Label)
	assert.Equal(t, "Jawaban B", result[0].Opsi[1].Isi)
	require.NoError(t, parser.ValidateParsedSoal(result))
}

func TestParseMarkersFromContent_EnumeratedOptionImageOnNextParagraph(t *testing.T) {
	t.Parallel()

	paragraphs := richParagraphs(
		"[Q:PG] Pilih gambar",
		"A. [A]",
		"[IMG]",
		"B. [B]",
		"[IMG]",
		"[ANS] A",
	)

	result, warnings, err := parser.ParseMarkersFromContent(paragraphs, buildDocxWithImages("image1.png", "image2.png"))
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, result, 1)
	require.Len(t, result[0].Opsi, 2)

	assert.Empty(t, result[0].Gambar)
	assert.Equal(t, "image1.png", result[0].Opsi[0].IsiContent.Blocks[0].Children[0].Src)
	assert.Equal(t, "image2.png", result[0].Opsi[1].IsiContent.Blocks[0].Children[0].Src)
	require.NoError(t, parser.ValidateParsedSoal(result))
}

func TestParseMarkersFromContent_StandaloneImageStaysQuestionImage(t *testing.T) {
	t.Parallel()

	paragraphs := richParagraphs(
		"[Q:PG] Perhatikan gambar",
		"[IMG]",
		"[A] Jawaban A",
		"[B] Jawaban B",
		"[ANS] A",
	)

	result, warnings, err := parser.ParseMarkersFromContent(paragraphs, buildDocxWithImages("image1.png"))
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, result, 1)

	assert.Equal(t, "image1.png", result[0].Gambar)
	assert.Empty(t, result[0].Opsi[0].IsiContent.Blocks[0].Children[0].Src)
}

func TestParseMarkersFromContent_OptionImageExplicitMissing(t *testing.T) {
	t.Parallel()

	paragraphs := richParagraphs(
		"[Q:PG] Pilih gambar yang benar",
		"[A] [IMG] image2.jpg Gambar pilihan",
		"[B] Jawaban biasa",
		"[ANS] A",
	)

	_, _, err := parser.ParseMarkersFromContent(paragraphs, buildDocxWithImages("image1.png"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), `gambar "image2.jpg" tidak ditemukan`)
}

func richParagraphs(items ...string) []content.RichContent {
	out := make([]content.RichContent, 0, len(items))
	for _, item := range items {
		out = append(out, content.RichContent{
			Version: 1,
			Blocks: []content.Block{
				{
					Type:     "paragraph",
					Children: []content.Inline{{Type: "text", Text: item}},
				},
			},
		})
	}
	return out
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
