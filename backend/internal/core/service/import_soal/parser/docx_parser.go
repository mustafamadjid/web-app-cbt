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

	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
)

// --- Relationships parsing (word/_rels/document.xml.rels) ---

type relationships struct {
	XMLName       xml.Name       `xml:"Relationships"`
	Relationships []relationship `xml:"Relationship"`
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
	contents, _, err := ExtractParagraphContents(data)
	if err != nil {
		return nil, err
	}
	paragraphs := make([]string, 0, len(contents))
	for _, item := range contents {
		text := strings.TrimSpace(item.PlainText())
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
	contentParagraphs := make([]content.RichContent, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		trimmed := strings.TrimSpace(paragraph)
		if trimmed == "" {
			continue
		}
		contentParagraphs = append(contentParagraphs, content.RichContent{
			Version: 1,
			Blocks: []content.Block{
				{
					Type: "paragraph",
					Children: []content.Inline{
						{Type: "text", Text: trimmed},
					},
				},
			},
		})
	}
	result, _, err := ParseMarkersFromContent(contentParagraphs, docxDataForAutoImg)
	return result, err
}

func parseWeightValue(raw string) (float64, error) {
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

// optionMarkerInfo checks whether a line starts with an option marker.
// It accepts the canonical "[A]" format and common exported-list variants
// like "A. [A]" where Word keeps both the visual list label and marker text.
func optionMarkerInfo(text string) (string, int, bool) {
	trimmed := strings.TrimLeft(text, " \t")
	leadingLen := len([]rune(text)) - len([]rune(trimmed))
	markerIdx := strings.Index(trimmed, "[")
	if markerIdx < 0 {
		return "", 0, false
	}

	prefix := strings.TrimSpace(trimmed[:markerIdx])
	marker := strings.ToUpper(trimmed[markerIdx:])
	if len(marker) < len("[A]") || marker[0] != '[' || marker[2] != ']' {
		return "", 0, false
	}

	label := string(marker[1])
	if !strings.Contains("ABCDE", label) {
		return "", 0, false
	}

	if prefix != "" {
		allowed := map[string]bool{
			label:       true,
			label + ".": true,
			label + ")": true,
			label + ":": true,
			label + "-": true,
		}
		if !allowed[strings.ToUpper(prefix)] {
			return "", 0, false
		}
	}

	prefixLen := leadingLen + len([]rune(trimmed[:markerIdx])) + len("[A]")
	return label, prefixLen, true
}
