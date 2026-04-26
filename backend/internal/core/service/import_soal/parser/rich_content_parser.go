package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
)

type xmlNode struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Nodes   []xmlNode  `xml:",any"`
	Text    string     `xml:",chardata"`
}

func ExtractParagraphContents(data []byte) ([]content.RichContent, []string, error) {
	documentBytes, err := readDocumentXML(data)
	if err != nil {
		return nil, nil, err
	}

	var root xmlNode
	if err := xml.Unmarshal(documentBytes, &root); err != nil {
		return nil, nil, fmt.Errorf("unmarshal document.xml: %w", err)
	}

	bodyNode := findFirstNode(root, "body")
	if bodyNode == nil {
		return nil, nil, fmt.Errorf("document body not found in docx")
	}

	paragraphs := make([]content.RichContent, 0)
	warnings := make([]string, 0)
	for _, node := range bodyNode.Nodes {
		if node.XMLName.Local != "p" {
			continue
		}

		rich, paragraphWarnings := paragraphNodeToContent(node)
		if !rich.Empty() {
			paragraphs = append(paragraphs, rich)
		}
		warnings = append(warnings, paragraphWarnings...)
	}

	return paragraphs, warnings, nil
}

func ParseMarkersFromContent(paragraphs []content.RichContent, docxDataForAutoImg []byte) ([]importsoal.ParsedSoal, []string, error) {
	var result []importsoal.ParsedSoal
	var current *importsoal.ParsedSoal
	var warnings []string

	imgQueue, err := extractImageNameQueue(docxDataForAutoImg)
	if err != nil {
		return nil, nil, fmt.Errorf("auto IMG mapping failed: %w", err)
	}
	imgIdx := 0
	nextImage := func() (string, bool) {
		if imgIdx >= len(imgQueue) {
			return "", false
		}
		name := imgQueue[imgIdx]
		imgIdx++
		return name, true
	}

	flushCurrent := func() {
		if current != nil {
			if current.Pertanyaan == "" {
				current.Pertanyaan = current.PertanyaanContent.PlainText()
			}
			for i := range current.Opsi {
				if current.Opsi[i].Isi == "" {
					current.Opsi[i].Isi = current.Opsi[i].IsiContent.PlainText()
				}
			}
			result = append(result, *current)
			current = nil
		}
	}

	type fillMode int
	const (
		modeQuestion fillMode = iota
		modeOption
		modeAns
		modeImg
		modeW
	)

	var mode fillMode = modeQuestion
	var lastOptLabel string

	appendContent := func(target *content.RichContent, block content.RichContent) {
		if target.Version == 0 {
			target.Version = 1
		}
		target.Blocks = append(target.Blocks, block.Blocks...)
	}

	appendQuestion := func(block content.RichContent) {
		if current == nil || block.Empty() {
			return
		}
		appendContent(&current.PertanyaanContent, block)
		current.Pertanyaan = current.PertanyaanContent.PlainText()
	}

	appendOption := func(label string, block content.RichContent) error {
		if block.Empty() {
			return nil
		}
		for i := range current.Opsi {
			if current.Opsi[i].Label == label {
				appendContent(&current.Opsi[i].IsiContent, block)
				current.Opsi[i].Isi = current.Opsi[i].IsiContent.PlainText()
				return nil
			}
		}
		current.Opsi = append(current.Opsi, importsoal.ParsedOpsi{
			Label:      label,
			Isi:        block.PlainText(),
			IsiContent: block,
		})
		return nil
	}

	parseWeight := func(raw string) (float64, error) {
		return parseWeightValue(raw)
	}

	markCorrect := func() {
		if current == nil || current.TipeSoal != "pilihan_ganda" {
			return
		}
		ans := strings.ToUpper(strings.TrimSpace(current.KunciJawaban))
		if ans == "" {
			return
		}
		for i := range current.Opsi {
			current.Opsi[i].IsBenar = (strings.ToUpper(current.Opsi[i].Label) == ans)
		}
	}

	for i, paragraph := range paragraphs {
		trimmed := strings.TrimSpace(paragraph.PlainText())
		if trimmed == "" {
			continue
		}
		upper := strings.ToUpper(trimmed)

		switch {
		case strings.HasPrefix(upper, "[Q:PG]"):
			flushCurrent()
			textContent := trimContentPrefix(paragraph, len("[Q:PG]"))
			current = &importsoal.ParsedSoal{
				Pertanyaan:        textContent.PlainText(),
				PertanyaanContent: textContent,
				TipeSoal:          "pilihan_ganda",
				BobotSoal:         1,
				NoUrutSoal:        len(result) + 1,
			}
			mode = modeQuestion
			lastOptLabel = ""

		case strings.HasPrefix(upper, "[Q:ESSAY]"):
			flushCurrent()
			textContent := trimContentPrefix(paragraph, len("[Q:ESSAY]"))
			current = &importsoal.ParsedSoal{
				Pertanyaan:        textContent.PlainText(),
				PertanyaanContent: textContent,
				TipeSoal:          "essay",
				BobotSoal:         1,
				NoUrutSoal:        len(result) + 1,
			}
			mode = modeQuestion
			lastOptLabel = ""

		case isOptionMarker(upper):
			if current == nil {
				return nil, warnings, fmt.Errorf("baris %d: opsi ditemukan tanpa soal aktif", i+1)
			}
			label := strings.ToUpper(string(trimmed[1]))
			textContent := trimContentPrefix(paragraph, len("[A]"))

			mode = modeOption
			lastOptLabel = label
			if err := appendOption(label, textContent); err != nil {
				return nil, warnings, fmt.Errorf("baris %d: %w", i+1, err)
			}

		case strings.HasPrefix(upper, "[ANS]"):
			if current == nil {
				return nil, warnings, fmt.Errorf("baris %d: [ANS] ditemukan tanpa soal aktif", i+1)
			}
			mode = modeAns
			raw := strings.TrimSpace(trimContentPrefix(paragraph, len("[ANS]")).PlainText())
			if raw != "" {
				current.KunciJawaban = raw
				markCorrect()
			}

		case strings.HasPrefix(upper, "[IMG]"):
			if current == nil {
				return nil, warnings, fmt.Errorf("baris %d: [IMG] ditemukan tanpa soal aktif", i+1)
			}
			mode = modeImg
			raw := strings.TrimSpace(trimContentPrefix(paragraph, len("[IMG]")).PlainText())
			if raw != "" {
				current.Gambar = raw
				break
			}
			if current.Gambar == "" {
				if name, ok := nextImage(); ok {
					current.Gambar = name
				} else {
					return nil, warnings, fmt.Errorf("baris %d: [IMG] tetapi tidak ada gambar tersisa di dokumen", i+1)
				}
			}

		case strings.HasPrefix(upper, "[W]"):
			if current == nil {
				return nil, warnings, fmt.Errorf("baris %d: [W] ditemukan tanpa soal aktif", i+1)
			}
			mode = modeW
			raw := strings.TrimSpace(trimContentPrefix(paragraph, len("[W]")).PlainText())
			if raw != "" {
				bobot, err := parseWeight(raw)
				if err != nil {
					return nil, warnings, fmt.Errorf("baris %d: %w", i+1, err)
				}
				current.BobotSoal = bobot
			}

		default:
			if current == nil {
				continue
			}
			switch mode {
			case modeQuestion:
				appendQuestion(paragraph)
			case modeOption:
				if lastOptLabel == "" {
					appendQuestion(paragraph)
					continue
				}
				if err := appendOption(lastOptLabel, paragraph); err != nil {
					return nil, warnings, fmt.Errorf("baris %d: %w", i+1, err)
				}
			case modeAns:
				raw := strings.TrimSpace(paragraph.PlainText())
				if current.KunciJawaban == "" {
					current.KunciJawaban = raw
				} else {
					current.KunciJawaban += "\n" + raw
				}
				markCorrect()
			case modeImg:
				if current.Gambar == "" {
					current.Gambar = strings.TrimSpace(paragraph.PlainText())
				}
			case modeW:
				bobot, err := parseWeight(strings.TrimSpace(paragraph.PlainText()))
				if err != nil {
					return nil, warnings, fmt.Errorf("baris %d: %w", i+1, err)
				}
				current.BobotSoal = bobot
			}
		}
	}

	flushCurrent()
	return result, warnings, nil
}

func paragraphNodeToContent(node xmlNode) (content.RichContent, []string) {
	result := content.New()
	block := content.Block{Type: "paragraph"}
	warnings := make([]string, 0)

	appendText := func(text string, marks []content.Mark) {
		if text == "" {
			return
		}
		text = strings.ReplaceAll(text, "\u00a0", " ")
		block.Children = append(block.Children, content.Inline{
			Type:  "text",
			Text:  text,
			Marks: cloneMarks(marks),
		})
	}

	var visitNode func(current xmlNode, inheritedMarks []content.Mark)
	visitNode = func(current xmlNode, inheritedMarks []content.Mark) {
		switch current.XMLName.Local {
		case "r":
			runMarks := mergeMarks(inheritedMarks, marksFromRun(current))
			for _, child := range current.Nodes {
				visitNode(child, runMarks)
			}
		case "t":
			appendText(current.Text, inheritedMarks)
		case "tab":
			appendText("\t", inheritedMarks)
		case "br":
			appendText("\n", inheritedMarks)
		case "oMath", "oMathPara":
			latex, warning := ommlNodeToLatex(current)
			if warning != "" {
				warnings = append(warnings, warning)
			}
			if strings.TrimSpace(latex) != "" {
				display := "inline"
				if current.XMLName.Local == "oMathPara" {
					display = "block"
				}
				block.Children = append(block.Children, content.Inline{
					Type:    "math",
					Latex:   latex,
					Display: display,
				})
			}
		default:
			for _, child := range current.Nodes {
				visitNode(child, inheritedMarks)
			}
		}
	}

	for _, child := range node.Nodes {
		visitNode(child, nil)
	}

	if len(block.Children) > 0 {
		result.Blocks = append(result.Blocks, block)
	}
	return result, warnings
}

func ommlNodeToLatex(node xmlNode) (string, string) {
	switch node.XMLName.Local {
	case "oMathPara":
		var parts []string
		var warnings []string
		for _, child := range node.Nodes {
			latex, warning := ommlNodeToLatex(child)
			if latex != "" {
				parts = append(parts, latex)
			}
			if warning != "" {
				warnings = append(warnings, warning)
			}
		}
		return strings.TrimSpace(strings.Join(parts, " ")), strings.Join(warnings, "; ")
	case "oMath":
		var parts []string
		var warnings []string
		for _, child := range node.Nodes {
			latex, warning := ommlNodeToLatex(child)
			if latex != "" {
				parts = append(parts, latex)
			}
			if warning != "" {
				warnings = append(warnings, warning)
			}
		}
		return strings.Join(parts, ""), strings.Join(warnings, "; ")
	case "r":
		return gatherOMMLText(node), ""
	case "t":
		return currentTextValue(node), ""
	case "sSup":
		base := ommlChildLatex(node, "e")
		sup := ommlChildLatex(node, "sup")
		if base == "" || sup == "" {
			return gatherOMMLText(node), "unsupported OMML superscript fallback ke raw text"
		}
		return fmt.Sprintf("{%s}^{%s}", base, sup), ""
	case "sSub":
		base := ommlChildLatex(node, "e")
		sub := ommlChildLatex(node, "sub")
		if base == "" || sub == "" {
			return gatherOMMLText(node), "unsupported OMML subscript fallback ke raw text"
		}
		return fmt.Sprintf("{%s}_{%s}", base, sub), ""
	case "sSubSup":
		base := ommlChildLatex(node, "e")
		sub := ommlChildLatex(node, "sub")
		sup := ommlChildLatex(node, "sup")
		if base == "" {
			return gatherOMMLText(node), "unsupported OMML subscript/superscript fallback ke raw text"
		}
		return fmt.Sprintf("{%s}_{%s}^{%s}", base, sub, sup), ""
	case "f":
		num := ommlChildLatex(node, "num")
		den := ommlChildLatex(node, "den")
		if num == "" || den == "" {
			return gatherOMMLText(node), "unsupported OMML fraction fallback ke raw text"
		}
		return fmt.Sprintf("\\frac{%s}{%s}", num, den), ""
	case "rad":
		deg := ommlChildLatex(node, "deg")
		base := ommlChildLatex(node, "e")
		if base == "" {
			return gatherOMMLText(node), "unsupported OMML radical fallback ke raw text"
		}
		if deg != "" {
			return fmt.Sprintf("\\sqrt[%s]{%s}", deg, base), ""
		}
		return fmt.Sprintf("\\sqrt{%s}", base), ""
	case "d":
		open, close := "(", ")"
		if beg := findFirstNode(node, "begChr"); beg != nil {
			if v := attrValue(*beg, "val"); v != "" {
				open = v
			}
		}
		if end := findFirstNode(node, "endChr"); end != nil {
			if v := attrValue(*end, "val"); v != "" {
				close = v
			}
		}
		body := ommlChildLatex(node, "e")
		return fmt.Sprintf("\\left%s %s \\right%s", open, body, close), ""
	case "nary":
		body := ommlChildLatex(node, "e")
		sub := ommlChildLatex(node, "sub")
		sup := ommlChildLatex(node, "sup")
		op := "\\sum"
		if chr := findFirstNode(node, "chr"); chr != nil {
			switch attrValue(*chr, "val") {
			case "∏":
				op = "\\prod"
			case "⋂":
				op = "\\bigcap"
			case "⋃":
				op = "\\bigcup"
			}
		}
		return op + optionalBound("_", sub) + optionalBound("^", sup) + "{" + body + "}", ""
	case "acc":
		base := ommlChildLatex(node, "e")
		if base == "" {
			return gatherOMMLText(node), "unsupported OMML accent fallback ke raw text"
		}
		cmd := "\\hat"
		if chr := findFirstNode(node, "chr"); chr != nil {
			switch attrValue(*chr, "val") {
			case "¯":
				cmd = "\\bar"
			case "→":
				cmd = "\\vec"
			case "˙":
				cmd = "\\dot"
			}
		}
		return fmt.Sprintf("%s{%s}", cmd, base), ""
	case "m":
		rows := make([]string, 0)
		for _, row := range node.Nodes {
			if row.XMLName.Local != "mr" {
				continue
			}
			cells := make([]string, 0)
			for _, cell := range row.Nodes {
				if cell.XMLName.Local != "e" {
					continue
				}
				cells = append(cells, ommlChildrenLatex(cell.Nodes))
			}
			rows = append(rows, strings.Join(cells, " & "))
		}
		if len(rows) == 0 {
			return gatherOMMLText(node), "unsupported OMML matrix fallback ke raw text"
		}
		return "\\begin{matrix}" + strings.Join(rows, ` \\ `) + "\\end{matrix}", ""
	case "func":
		name := ommlChildLatex(node, "fName")
		arg := ommlChildLatex(node, "e")
		return name + "{" + arg + "}", ""
	default:
		text := gatherOMMLText(node)
		if strings.TrimSpace(text) == "" {
			return "", ""
		}
		return text, fmt.Sprintf("unsupported OMML node %q fallback ke raw text", node.XMLName.Local)
	}
}

func readDocumentXML(data []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open document.xml: %w", err)
		}
		defer rc.Close()
		xmlBytes, err := io.ReadAll(rc)
		if err != nil {
			return nil, fmt.Errorf("read document.xml: %w", err)
		}
		return xmlBytes, nil
	}

	return nil, fmt.Errorf("word/document.xml not found in docx")
}

func findFirstNode(node xmlNode, local string) *xmlNode {
	if node.XMLName.Local == local {
		return &node
	}
	for _, child := range node.Nodes {
		if child.XMLName.Local == local {
			return &child
		}
		if nested := findFirstNode(child, local); nested != nil {
			return nested
		}
	}
	return nil
}

func currentTextValue(node xmlNode) string {
	var sb strings.Builder
	var visit func(xmlNode)
	visit = func(item xmlNode) {
		if item.XMLName.Local == "t" && item.Text != "" {
			sb.WriteString(item.Text)
		}
		if item.Text != "" && item.XMLName.Local == "" {
			sb.WriteString(item.Text)
		}
		for _, child := range item.Nodes {
			visit(child)
		}
	}
	visit(node)
	return sb.String()
}

func gatherOMMLText(node xmlNode) string {
	var sb strings.Builder
	var visit func(xmlNode)
	visit = func(item xmlNode) {
		if item.Text != "" {
			sb.WriteString(item.Text)
		}
		for _, child := range item.Nodes {
			visit(child)
		}
	}
	visit(node)
	return strings.TrimSpace(sb.String())
}

func ommlChildLatex(node xmlNode, local string) string {
	for _, child := range node.Nodes {
		if child.XMLName.Local == local {
			return ommlChildrenLatex(child.Nodes)
		}
	}
	return ""
}

func ommlChildrenLatex(nodes []xmlNode) string {
	parts := make([]string, 0, len(nodes))
	for _, child := range nodes {
		latex, _ := ommlNodeToLatex(child)
		if latex != "" {
			parts = append(parts, latex)
		}
	}
	return strings.Join(parts, "")
}

func attrValue(node xmlNode, local string) string {
	for _, attr := range node.Attrs {
		if attr.Name.Local == local {
			return attr.Value
		}
	}
	return ""
}

func cloneMarks(marks []content.Mark) []content.Mark {
	if len(marks) == 0 {
		return nil
	}
	out := make([]content.Mark, len(marks))
	copy(out, marks)
	return out
}

func mergeMarks(base, extra []content.Mark) []content.Mark {
	out := cloneMarks(base)
	for _, mark := range extra {
		exists := false
		for _, current := range out {
			if current == mark {
				exists = true
				break
			}
		}
		if !exists {
			out = append(out, mark)
		}
	}
	return out
}

func marksFromRun(node xmlNode) []content.Mark {
	props := findFirstNode(node, "rPr")
	if props == nil {
		return nil
	}
	vertAlign := findFirstNode(*props, "vertAlign")
	if vertAlign == nil {
		return nil
	}
	switch attrValue(*vertAlign, "val") {
	case "superscript":
		return []content.Mark{content.MarkSup}
	case "subscript":
		return []content.Mark{content.MarkSub}
	default:
		return nil
	}
}

func trimContentPrefix(value content.RichContent, prefixLen int) content.RichContent {
	if prefixLen <= 0 || len(value.Blocks) == 0 {
		return value
	}

	remaining := prefixLen
	result := content.New()
	for _, block := range value.Blocks {
		nextBlock := content.Block{Type: block.Type}
		for _, child := range block.Children {
			switch child.Type {
			case "text":
				if remaining > 0 {
					runes := []rune(child.Text)
					if len(runes) <= remaining {
						remaining -= len(runes)
						continue
					}
					child.Text = string(runes[remaining:])
					remaining = 0
				}
				if child.Text != "" {
					nextBlock.Children = append(nextBlock.Children, child)
				}
			default:
				nextBlock.Children = append(nextBlock.Children, child)
			}
		}
		if len(nextBlock.Children) > 0 {
			result.Blocks = append(result.Blocks, nextBlock)
		}
	}

	return trimLeadingWhitespaceContent(result)
}

func trimLeadingWhitespaceContent(value content.RichContent) content.RichContent {
	for blockIdx := range value.Blocks {
		for childIdx := range value.Blocks[blockIdx].Children {
			child := &value.Blocks[blockIdx].Children[childIdx]
			if child.Type != "text" {
				return value
			}
			trimmed := strings.TrimLeft(child.Text, " \t")
			if trimmed == "" {
				child.Text = ""
				continue
			}
			child.Text = trimmed
			return compactContent(value)
		}
	}
	return compactContent(value)
}

func compactContent(value content.RichContent) content.RichContent {
	result := content.New()
	for _, block := range value.Blocks {
		nextBlock := content.Block{Type: block.Type}
		for _, child := range block.Children {
			if child.Type == "text" && child.Text == "" {
				continue
			}
			nextBlock.Children = append(nextBlock.Children, child)
		}
		if len(nextBlock.Children) > 0 {
			result.Blocks = append(result.Blocks, nextBlock)
		}
	}
	return result
}

func optionalBound(prefix, value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return prefix + "{" + value + "}"
}
