package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"path"
	"strconv"
	"strings"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
)

// --- OOXML structs for parsing word/document.xml ---

type document struct {
	Body body `xml:"body"`
}

type body struct {
	Paragraphs []paragraph `xml:"p"`
}

type paragraph struct {
	Runs []run `xml:"r"`
}

type run struct {
	Text runText `xml:"t"`
	// Drawing exists in real OOXML but we don't rely on unmarshalling it here;
	// we detect images via token scan for <a:blip r:embed="...">.
}

type runText struct {
	Value string `xml:",chardata"`
}

// --- Relationships parsing (word/_rels/document.xml.rels) ---

type relationships struct {
	XMLName       xml.Name       `xml:"Relationships"`
	Relationships []relationship  `xml:"Relationship"`
}

type relationship struct {
	Id     string `xml:"Id,attr"`
	Type   string `xml:"Type,attr"`
	Target string `xml:"Target,attr"`
}

// ExtractParagraphs opens a .docx byte slice (ZIP), finds word/document.xml,
// and returns the text content of every <w:p> paragraph.
// Note: this is still "text-only"; image mapping is handled by functions below.
func ExtractParagraphs(data []byte) ([]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	var docFile *zip.File
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			docFile = f
			break
		}
	}
	if docFile == nil {
		return nil, fmt.Errorf("word/document.xml not found in docx")
	}

	rc, err := docFile.Open()
	if err != nil {
		return nil, fmt.Errorf("open document.xml: %w", err)
	}
	defer rc.Close()

	xmlBytes, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read document.xml: %w", err)
	}

	var doc document
	if err := xml.Unmarshal(xmlBytes, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal document.xml: %w", err)
	}

	paragraphs := make([]string, 0, len(doc.Body.Paragraphs))
	for _, p := range doc.Body.Paragraphs {
		var sb strings.Builder
		for _, r := range p.Runs {
			sb.WriteString(r.Text.Value)
		}
		text := strings.TrimSpace(sb.String())
		if text != "" {
			paragraphs = append(paragraphs, text)
		}
	}

	return paragraphs, nil
}

func ExtractImageFiles(data []byte) (map[string][]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	images := make(map[string][]byte)
	for _, f := range zr.File {
		if !strings.HasPrefix(f.Name, "word/media/") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open media file %s: %w", f.Name, err)
		}
		b, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read media file %s: %w", f.Name, err)
		}
		name := path.Base(f.Name) // image1.png
		images[name] = b
	}
	return images, nil
}

// --- NEW: read relationship map rId -> "word/media/xxx.png" ---

func extractDocRelsMap(data []byte) (map[string]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	var relFile *zip.File
	for _, f := range zr.File {
		if f.Name == "word/_rels/document.xml.rels" {
			relFile = f
			break
		}
	}
	if relFile == nil {
		return nil, fmt.Errorf("word/_rels/document.xml.rels not found in docx")
	}

	rc, err := relFile.Open()
	if err != nil {
		return nil, fmt.Errorf("open document.xml.rels: %w", err)
	}
	defer rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read document.xml.rels: %w", err)
	}

	var rels relationships
	if err := xml.Unmarshal(b, &rels); err != nil {
		return nil, fmt.Errorf("unmarshal document.xml.rels: %w", err)
	}

	// Only keep image relationships
	// Type for images commonly ends with ".../relationships/image"
	out := make(map[string]string, len(rels.Relationships))
	for _, r := range rels.Relationships {
		if r.Id == "" || r.Target == "" {
			continue
		}
		if !strings.Contains(strings.ToLower(r.Type), "/image") {
			continue
		}

		// Target is typically "media/image1.png"
		target := strings.TrimSpace(r.Target)
		target = strings.TrimPrefix(target, "/")
		target = strings.TrimPrefix(target, "../")
		if !strings.HasPrefix(target, "word/") {
			target = "word/" + target
		}
		out[r.Id] = target // rId5 -> word/media/image1.png
	}

	return out, nil
}

// --- NEW: scan document.xml to get image order by r:embed (rId) ---

func extractImageRidOrderFromDocumentXML(data []byte) ([]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	var docFile *zip.File
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			docFile = f
			break
		}
	}
	if docFile == nil {
		return nil, fmt.Errorf("word/document.xml not found in docx")
	}

	rc, err := docFile.Open()
	if err != nil {
		return nil, fmt.Errorf("open document.xml: %w", err)
	}
	defer rc.Close()

	dec := xml.NewDecoder(rc)

	var rids []string
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decode document.xml: %w", err)
		}

		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}

		// In OOXML, images are referenced via:
		// <a:blip r:embed="rIdX" ... />
		// We detect any element with local name "blip" and an attribute local name "embed"
		// in the relationships namespace or any prefix.
		if strings.EqualFold(se.Name.Local, "blip") {
			for _, a := range se.Attr {
				if strings.EqualFold(a.Name.Local, "embed") && strings.TrimSpace(a.Value) != "" {
					rids = append(rids, strings.TrimSpace(a.Value))
					break
				}
			}
		}
	}

	return rids, nil
}

// --- NEW: resolve image order to actual media basenames (image1.png, etc.) ---

func extractImageNameQueue(data []byte) ([]string, error) {
	relMap, err := extractDocRelsMap(data)
	if err != nil {
		return nil, err
	}

	rids, err := extractImageRidOrderFromDocumentXML(data)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, rid := range rids {
		target, ok := relMap[rid]
		if !ok {
			continue
		}
		// target is "word/media/image1.png"
		names = append(names, path.Base(target))
	}

	return names, nil
}

// --- Marker parsing (same as your original, but [IMG] auto-mapped) ---

// ParseMarkers parses paragraphs and extracts structured exam questions based on marker tags.
// Supported markers:
//
//	[Q:PG]    — start a multiple-choice question
//	[Q:ESSAY] — start an essay question
//	[A]-[E]   — multiple-choice option
//	[ANS]     — answer key
//	[IMG]     — (NOW) optional image filename; if empty, auto-attach next inserted image
//	[W]       — weight / bobot

func ParseMarkers(paragraphs []string, docxDataForAutoImg []byte) ([]importsoal.ParsedSoal, error) {
	var result []importsoal.ParsedSoal
	var current *importsoal.ParsedSoal

	imgQueue, err := extractImageNameQueue(docxDataForAutoImg)
	if err != nil {
		return nil, fmt.Errorf("auto IMG mapping failed: %w", err)
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

	appendQuestion := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" || current == nil {
			return
		}
		if current.Pertanyaan == "" {
			current.Pertanyaan = s
		} else {
			current.Pertanyaan += "\n" + s
		}
	}

	appendOption := func(label, s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		// append ke opsi terakhir dengan label yang sama (multi-line)
		for i := range current.Opsi {
			if current.Opsi[i].Label == label {
				if current.Opsi[i].Isi == "" {
					current.Opsi[i].Isi = s
				} else {
					current.Opsi[i].Isi += "\n" + s
				}
				return nil
			}
		}
		current.Opsi = append(current.Opsi, importsoal.ParsedOpsi{
			Label: label,
			Isi:   s,
		})
		return nil
	}

	parseWeight := func(raw string) (float64, error) {
		raw = strings.TrimSpace(raw)
		bobot, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return 0, fmt.Errorf("bobot bukan angka: %q", raw)
		}
		if math.IsNaN(bobot) || math.IsInf(bobot, 0) {
			return 0, fmt.Errorf("bobot bukan angka: %q", raw)
		}
		if bobot <= 0 {
			return 0, fmt.Errorf("bobot harus lebih dari 0")
		}
		return bobot, nil
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

	for i, line := range paragraphs {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		upper := strings.ToUpper(trimmed)

		switch {
		case strings.HasPrefix(upper, "[Q:PG]"):
			flushCurrent()
			teks := strings.TrimSpace(trimmed[len("[Q:PG]"):])
			current = &importsoal.ParsedSoal{
				Pertanyaan: teks, // bisa kosong; nanti diisi paragraf berikutnya
				TipeSoal:   "pilihan_ganda",
				BobotSoal:  1,
				NoUrutSoal: len(result) + 1,
			}
			mode = modeQuestion
			lastOptLabel = ""

		case strings.HasPrefix(upper, "[Q:ESSAY]"):
			flushCurrent()
			teks := strings.TrimSpace(trimmed[len("[Q:ESSAY]"):])
			current = &importsoal.ParsedSoal{
				Pertanyaan: teks,
				TipeSoal:   "essay",
				BobotSoal:  1,
				NoUrutSoal: len(result) + 1,
			}
			mode = modeQuestion
			lastOptLabel = ""

		case isOptionMarker(upper):
			if current == nil {
				return nil, fmt.Errorf("baris %d: opsi ditemukan tanpa soal aktif", i+1)
			}
			label := strings.ToUpper(string(trimmed[1]))
			text := ""
			if len(trimmed) > 3 {
				text = strings.TrimSpace(trimmed[3:]) // setelah "]"
			}
			// hapus leading spasi kalau ada
			text = strings.TrimLeft(text, " \t")

			mode = modeOption
			lastOptLabel = label
			if err := appendOption(label, text); err != nil {
				return nil, fmt.Errorf("baris %d: %w", i+1, err)
			}

		case strings.HasPrefix(upper, "[ANS]"):
			if current == nil {
				return nil, fmt.Errorf("baris %d: [ANS] ditemukan tanpa soal aktif", i+1)
			}
			mode = modeAns
			raw := strings.TrimSpace(trimmed[len("[ANS]"):])
			if raw != "" {
				current.KunciJawaban = raw
				markCorrect()
			}

		case strings.HasPrefix(upper, "[IMG]"):
			if current == nil {
				return nil, fmt.Errorf("baris %d: [IMG] ditemukan tanpa soal aktif", i+1)
			}
			mode = modeImg
			raw := strings.TrimSpace(trimmed[len("[IMG]"):])
			if raw != "" {
				current.Gambar = raw
				break
			}
			// isi otomatis kalau marker kosong
			if current.Gambar == "" {
				if name, ok := nextImage(); ok {
					current.Gambar = name
				} else {
					return nil, fmt.Errorf("baris %d: [IMG] tetapi tidak ada gambar tersisa di dokumen", i+1)
				}
			}

		case strings.HasPrefix(upper, "[W]"):
			if current == nil {
				return nil, fmt.Errorf("baris %d: [W] ditemukan tanpa soal aktif", i+1)
			}
			mode = modeW
			raw := strings.TrimSpace(trimmed[len("[W]"):])
			if raw != "" {
				bobot, err := parseWeight(raw)
				if err != nil {
					return nil, fmt.Errorf("baris %d: %w", i+1, err)
				}
				current.BobotSoal = bobot
			}

		default:
			// Isi berdasarkan mode terakhir
			if current == nil {
				continue
			}
			switch mode {
			case modeQuestion:
				appendQuestion(trimmed)
			case modeOption:
				if lastOptLabel == "" {
					appendQuestion(trimmed)
					continue
				}
				if err := appendOption(lastOptLabel, trimmed); err != nil {
					return nil, fmt.Errorf("baris %d: %w", i+1, err)
				}
			case modeAns:
				if current.KunciJawaban == "" {
					current.KunciJawaban = trimmed
				} else {
					current.KunciJawaban += "\n" + trimmed
				}
				markCorrect()
			case modeImg:
				// kalau user nulis [IMG] lalu nama file di paragraf berikutnya
				if current.Gambar == "" {
					current.Gambar = trimmed
				}
			case modeW:
				// kalau user nulis [W] lalu angka di paragraf berikutnya
				bobot, err := parseWeight(trimmed)
				if err != nil {
					return nil, fmt.Errorf("baris %d: %w", i+1, err)
				}
				current.BobotSoal = bobot
			}
		}
	}

	flushCurrent()
	return result, nil
}

// ValidateParsedSoal validates all parsed questions.
func ValidateParsedSoal(soalList []importsoal.ParsedSoal) error {
	if len(soalList) == 0 {
		return fmt.Errorf("tidak ada soal yang ditemukan")
	}

	for i, s := range soalList {
		num := i + 1
		if strings.TrimSpace(s.Pertanyaan) == "" {
			return fmt.Errorf("soal ke-%d: pertanyaan kosong", num)
		}
		if s.BobotSoal <= 0 {
			return fmt.Errorf("soal ke-%d: bobot harus lebih dari 0", num)
		}

		switch s.TipeSoal {
		case "pilihan_ganda":
			if len(s.Opsi) < 2 {
				return fmt.Errorf("soal ke-%d: pilihan ganda harus punya minimal 2 opsi", num)
			}
			if s.KunciJawaban == "" {
				return fmt.Errorf("soal ke-%d: kunci jawaban kosong", num)
			}
			jawaban := strings.ToUpper(s.KunciJawaban)
			found := false
			for _, o := range s.Opsi {
				if o.Label == jawaban {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("soal ke-%d: kunci jawaban %q tidak cocok dengan opsi", num, s.KunciJawaban)
			}
		case "essay":
			// ok
		default:
			return fmt.Errorf("soal ke-%d: tipe soal tidak valid: %q", num, s.TipeSoal)
		}
	}
	return nil
}

// isOptionMarker checks if a line starts with [A], [B], [C], [D], or [E]
func isOptionMarker(upper string) bool {
	for _, label := range []string{"[A]", "[B]", "[C]", "[D]", "[E]"} {
		if strings.HasPrefix(upper, label) {
			return true
		}
	}
	return false
}