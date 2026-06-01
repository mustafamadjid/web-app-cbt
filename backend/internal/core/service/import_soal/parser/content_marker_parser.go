package parser

import (
	"fmt"
	"strings"

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
)

type markerFillMode int

const (
	modeQuestion markerFillMode = iota
	modeOption
	modeAns
	modeImg
	modeW
)

type contentMarkerParser struct {
	result       []importsoal.ParsedSoal
	current      *importsoal.ParsedSoal
	warnings     []string
	images       *imageNameQueue
	mode         markerFillMode
	lastOptLabel string
}

func ParseMarkersFromContent(paragraphs []content.RichContent, docxDataForAutoImg []byte) ([]importsoal.ParsedSoal, []string, error) {
	images, err := newImageNameQueue(docxDataForAutoImg)
	if err != nil {
		return nil, nil, fmt.Errorf("auto IMG mapping failed: %w", err)
	}

	parser := &contentMarkerParser{
		images: images,
		mode:   modeQuestion,
	}
	for i, paragraph := range paragraphs {
		if err := parser.parseParagraph(i+1, paragraph); err != nil {
			return nil, parser.warnings, err
		}
	}

	parser.flushCurrent()
	return parser.result, parser.warnings, nil
}

func (p *contentMarkerParser) parseParagraph(line int, paragraph content.RichContent) error {
	trimmed := strings.TrimSpace(paragraph.PlainText())
	if trimmed == "" {
		return nil
	}

	upper := strings.ToUpper(trimmed)
	optionLabel, markerPrefixLen, isOption := optionMarkerInfo(trimmed)

	switch {
	case strings.HasPrefix(upper, "[Q:PG]"):
		p.startQuestion(paragraph, len("[Q:PG]"), "pilihan_ganda")
	case strings.HasPrefix(upper, "[Q:ESSAY]"):
		p.startQuestion(paragraph, len("[Q:ESSAY]"), "essay")
	case isOption:
		return p.parseOptionMarker(line, paragraph, optionLabel, markerPrefixLen)
	case strings.HasPrefix(upper, "[ANS]"):
		return p.parseAnswerMarker(line, paragraph)
	case strings.HasPrefix(upper, "[IMG]"):
		return p.parseImageMarker(line, paragraph)
	case strings.HasPrefix(upper, "[W]"):
		return p.parseWeightMarker(line, paragraph)
	default:
		return p.parseContinuation(line, paragraph)
	}

	return nil
}

func (p *contentMarkerParser) startQuestion(paragraph content.RichContent, markerLen int, tipeSoal string) {
	p.flushCurrent()
	textContent := trimContentPrefix(paragraph, markerLen)
	p.current = &importsoal.ParsedSoal{
		Pertanyaan:        textContent.PlainText(),
		PertanyaanContent: textContent,
		TipeSoal:          tipeSoal,
		BobotSoal:         1,
		NoUrutSoal:        len(p.result) + 1,
	}
	p.mode = modeQuestion
	p.lastOptLabel = ""
}

func (p *contentMarkerParser) parseOptionMarker(line int, paragraph content.RichContent, label string, markerPrefixLen int) error {
	if p.current == nil {
		return fmt.Errorf("baris %d: opsi ditemukan tanpa soal aktif", line)
	}

	textContent := trimContentPrefix(paragraph, markerPrefixLen)
	textContent, err := p.parseOptionContent(textContent)
	if err != nil {
		return fmt.Errorf("baris %d: %w", line, err)
	}

	p.mode = modeOption
	p.lastOptLabel = label
	if err := p.appendOption(label, textContent); err != nil {
		return fmt.Errorf("baris %d: %w", line, err)
	}
	return nil
}

func (p *contentMarkerParser) parseAnswerMarker(line int, paragraph content.RichContent) error {
	if p.current == nil {
		return fmt.Errorf("baris %d: [ANS] ditemukan tanpa soal aktif", line)
	}

	p.mode = modeAns
	raw := strings.TrimSpace(trimContentPrefix(paragraph, len("[ANS]")).PlainText())
	if raw != "" {
		p.current.KunciJawaban = raw
		p.markCorrect()
	}
	return nil
}

func (p *contentMarkerParser) parseImageMarker(line int, paragraph content.RichContent) error {
	if p.current == nil {
		return fmt.Errorf("baris %d: [IMG] ditemukan tanpa soal aktif", line)
	}
	if p.mode == modeOption && p.lastOptLabel != "" {
		return p.appendImageToLastOption(line, paragraph)
	}

	p.mode = modeImg
	raw := strings.TrimSpace(trimContentPrefix(paragraph, len("[IMG]")).PlainText())
	if raw != "" {
		p.current.Gambar = raw
		return nil
	}
	if p.current.Gambar != "" {
		return nil
	}

	name, ok := p.images.next()
	if !ok {
		return fmt.Errorf("baris %d: [IMG] tetapi tidak ada gambar tersisa di dokumen", line)
	}
	p.current.Gambar = name
	return nil
}

func (p *contentMarkerParser) appendImageToLastOption(line int, paragraph content.RichContent) error {
	textContent, err := p.parseOptionContent(paragraph)
	if err != nil {
		return fmt.Errorf("baris %d: %w", line, err)
	}
	if err := p.appendOption(p.lastOptLabel, textContent); err != nil {
		return fmt.Errorf("baris %d: %w", line, err)
	}
	return nil
}

func (p *contentMarkerParser) parseWeightMarker(line int, paragraph content.RichContent) error {
	if p.current == nil {
		return fmt.Errorf("baris %d: [W] ditemukan tanpa soal aktif", line)
	}

	p.mode = modeW
	raw := strings.TrimSpace(trimContentPrefix(paragraph, len("[W]")).PlainText())
	if raw == "" {
		return nil
	}

	bobot, err := parseWeightValue(raw)
	if err != nil {
		return fmt.Errorf("baris %d: %w", line, err)
	}
	p.current.BobotSoal = bobot
	return nil
}

func (p *contentMarkerParser) parseContinuation(line int, paragraph content.RichContent) error {
	if p.current == nil {
		return nil
	}

	switch p.mode {
	case modeQuestion:
		p.appendQuestion(paragraph)
	case modeOption:
		return p.appendOptionContinuation(line, paragraph)
	case modeAns:
		p.appendAnswerContinuation(paragraph)
	case modeImg:
		if p.current.Gambar == "" {
			p.current.Gambar = strings.TrimSpace(paragraph.PlainText())
		}
	case modeW:
		return p.appendWeightContinuation(line, paragraph)
	}
	return nil
}

func (p *contentMarkerParser) appendOptionContinuation(line int, paragraph content.RichContent) error {
	if p.lastOptLabel == "" {
		p.appendQuestion(paragraph)
		return nil
	}
	if err := p.appendOption(p.lastOptLabel, paragraph); err != nil {
		return fmt.Errorf("baris %d: %w", line, err)
	}
	return nil
}

func (p *contentMarkerParser) appendAnswerContinuation(paragraph content.RichContent) {
	raw := strings.TrimSpace(paragraph.PlainText())
	if p.current.KunciJawaban == "" {
		p.current.KunciJawaban = raw
	} else {
		p.current.KunciJawaban += "\n" + raw
	}
	p.markCorrect()
}

func (p *contentMarkerParser) appendWeightContinuation(line int, paragraph content.RichContent) error {
	bobot, err := parseWeightValue(strings.TrimSpace(paragraph.PlainText()))
	if err != nil {
		return fmt.Errorf("baris %d: %w", line, err)
	}
	p.current.BobotSoal = bobot
	return nil
}

func (p *contentMarkerParser) flushCurrent() {
	if p.current == nil {
		return
	}

	if p.current.Pertanyaan == "" {
		p.current.Pertanyaan = p.current.PertanyaanContent.PlainText()
	}
	for i := range p.current.Opsi {
		if p.current.Opsi[i].Isi == "" {
			p.current.Opsi[i].Isi = p.current.Opsi[i].IsiContent.PlainText()
		}
	}
	p.result = append(p.result, *p.current)
	p.current = nil
}

func (p *contentMarkerParser) appendQuestion(block content.RichContent) {
	if p.current == nil || block.Empty() {
		return
	}
	appendContent(&p.current.PertanyaanContent, block)
	p.current.Pertanyaan = p.current.PertanyaanContent.PlainText()
}

func (p *contentMarkerParser) appendOption(label string, block content.RichContent) error {
	if block.Empty() {
		return nil
	}
	for i := range p.current.Opsi {
		if p.current.Opsi[i].Label == label {
			appendContent(&p.current.Opsi[i].IsiContent, block)
			p.current.Opsi[i].Isi = p.current.Opsi[i].IsiContent.PlainText()
			return nil
		}
	}
	p.current.Opsi = append(p.current.Opsi, importsoal.ParsedOpsi{
		Label:      label,
		Isi:        block.PlainText(),
		IsiContent: block,
	})
	return nil
}

func (p *contentMarkerParser) parseOptionContent(block content.RichContent) (content.RichContent, error) {
	block = trimLeadingWhitespaceContent(block)
	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(block.PlainText())), "[IMG]") {
		return block, nil
	}

	imageBody := trimContentPrefix(block, len("[IMG]"))
	imageBody = trimLeadingWhitespaceContent(imageBody)
	rawBody := strings.TrimSpace(imageBody.PlainText())

	src, captionContent, err := imageContentParts(
		rawBody,
		imageBody,
		p.images.next,
		p.images.exists,
		p.images.consumeIfNext,
	)
	if err != nil {
		return content.RichContent{}, err
	}

	imageContent := content.RichContent{
		Version: 1,
		Blocks: []content.Block{
			{
				Type: "paragraph",
				Children: []content.Inline{
					{Type: "image", Src: src, Alt: captionContent.PlainText()},
				},
			},
		},
	}
	appendContent(&imageContent, captionContent)
	return imageContent, nil
}

func (p *contentMarkerParser) markCorrect() {
	if p.current == nil || p.current.TipeSoal != "pilihan_ganda" {
		return
	}

	ans := strings.ToUpper(strings.TrimSpace(p.current.KunciJawaban))
	if ans == "" {
		return
	}
	for i := range p.current.Opsi {
		p.current.Opsi[i].IsBenar = strings.ToUpper(p.current.Opsi[i].Label) == ans
	}
}
