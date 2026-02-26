package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
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
	Text    runText `xml:"t"`
	Drawing drawing `xml:"drawing"`
}

type runText struct {
	Value string `xml:",chardata"`
}

// Minimal drawing detection — we only care whether a run contains as drawing.
type drawing struct {
	XMLName xml.Name
}

// ExtractParagraphs opens a .docx byte slice (ZIP), finds word/document.xml,
// and returns the text content of every <w:p> paragraph.
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

// ExtractImageFiles returns a map of image filename -> image bytes from the
// DOCX ZIP archive (word/media/*).
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
		// Use just the filename, not the full path
		parts := strings.Split(f.Name, "/")
		name := parts[len(parts)-1]
		images[name] = b
	}

	return images, nil
}

// ParseMarkers parses a list of paragraph strings and extracts structured
// exam questions based on marker tags.
//
// Supported markers:
//
//	[Q:PG]    — start a multiple-choice question
//	[Q:ESSAY] — start an essay question
//	[A]-[E]   — multiple-choice option
//	[ANS]     — answer key
//	[IMG]     — image filename
//	[W]       — weight / bobot
func ParseMarkers(paragraphs []string) ([]importsoal.ParsedSoal, error) {
	var result []importsoal.ParsedSoal
	var current *importsoal.ParsedSoal

	flushCurrent := func() {
		if current != nil {
			result = append(result, *current)
			current = nil
		}
	}

	for i, line := range paragraphs {
		upper := strings.ToUpper(strings.TrimSpace(line))

		switch {
		case strings.HasPrefix(upper, "[Q:PG]"):
			flushCurrent()
			teks := strings.TrimSpace(line[len("[Q:PG]"):])
			current = &importsoal.ParsedSoal{
				Pertanyaan: teks,
				TipeSoal:   "pilihan_ganda",
				BobotSoal:  1,
			}

		case strings.HasPrefix(upper, "[Q:ESSAY]"):
			flushCurrent()
			teks := strings.TrimSpace(line[len("[Q:ESSAY]"):])
			current = &importsoal.ParsedSoal{
				Pertanyaan: teks,
				TipeSoal:   "essay",
				BobotSoal:  1,
			}

		case isOptionMarker(upper):
			if current == nil {
				return nil, fmt.Errorf("baris %d: opsi ditemukan tanpa soal aktif", i+1)
			}
			label := strings.ToUpper(string(line[1]))
			isi := strings.TrimSpace(line[4:]) // skip "[X] "
			current.Opsi = append(current.Opsi, importsoal.ParsedOpsi{
				Label: label,
				Isi:   isi,
			})

		case strings.HasPrefix(upper, "[ANS]"):
			if current == nil {
				return nil, fmt.Errorf("baris %d: [ANS] ditemukan tanpa soal aktif", i+1)
			}
			jawaban := strings.TrimSpace(line[len("[ANS]"):])
			current.KunciJawaban = jawaban
			// Mark the correct option for PG
			if current.TipeSoal == "pilihan_ganda" {
				jawabanUpper := strings.ToUpper(jawaban)
				for idx := range current.Opsi {
					if current.Opsi[idx].Label == jawabanUpper {
						current.Opsi[idx].IsBenar = true
					}
				}
			}

		case strings.HasPrefix(upper, "[IMG]"):
			if current == nil {
				return nil, fmt.Errorf("baris %d: [IMG] ditemukan tanpa soal aktif", i+1)
			}
			current.Gambar = strings.TrimSpace(line[len("[IMG]"):])

		case strings.HasPrefix(upper, "[W]"):
			if current == nil {
				return nil, fmt.Errorf("baris %d: [W] ditemukan tanpa soal aktif", i+1)
			}
			var bobot int
			raw := strings.TrimSpace(line[len("[W]"):])
			if _, err := fmt.Sscanf(raw, "%d", &bobot); err != nil {
				return nil, fmt.Errorf("baris %d: bobot bukan angka: %q", i+1, raw)
			}
			if bobot <= 0 {
				return nil, fmt.Errorf("baris %d: bobot harus lebih dari 0", i+1)
			}
			current.BobotSoal = bobot

		default:
			// Baris tanpa marker — append ke pertanyaan soal aktif (multi-line pertanyaan)
			if current != nil && current.Pertanyaan != "" {
				current.Pertanyaan += "\n" + line
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
			// Essay hanya perlu pertanyaan dan bobot
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
