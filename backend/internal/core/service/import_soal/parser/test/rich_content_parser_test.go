package parser_test

import (
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
