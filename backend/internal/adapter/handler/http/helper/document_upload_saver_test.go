package httpx

import (
	"archive/zip"
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type inMemoryMultipartFile struct {
	*bytes.Reader
}

func (f *inMemoryMultipartFile) Close() error {
	return nil
}

func TestSaveDocumentRelative_AcceptsDocxOOXMLZip(t *testing.T) {
	docxBytes := buildOOXMLZip(t, map[string]string{
		"[Content_Types].xml": `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"></Types>`,
		"word/document.xml":   `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"></w:document>`,
	})

	store := DocumentStore{
		Dir:      t.TempDir(),
		Route:    "/uploads",
		MaxBytes: 5 << 20,
	}

	file := &inMemoryMultipartFile{Reader: bytes.NewReader(docxBytes)}
	fh := &multipart.FileHeader{
		Filename: "Surat Edaran UTS.docx",
		Size:     int64(len(docxBytes)),
	}

	relativePath, err := store.SaveDocumentRelative(file, fh)
	if err != nil {
		t.Fatalf("SaveDocumentRelative returned error: %v", err)
	}

	expectedFileName := time.Now().Format("20060102") + "_Surat Edaran UTS.docx"
	if relativePath != "/uploads/"+expectedFileName {
		t.Fatalf("expected relative path %q, got %q", "/uploads/"+expectedFileName, relativePath)
	}

	savedPath := filepath.Join(store.Dir, filepath.Base(relativePath))
	if _, err := os.Stat(savedPath); err != nil {
		t.Fatalf("saved file is missing: %v", err)
	}
}

func TestSaveDocumentRelative_UsesIncrementSuffixWhenFileExists(t *testing.T) {
	docxBytes := buildOOXMLZip(t, map[string]string{
		"[Content_Types].xml": `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"></Types>`,
		"word/document.xml":   `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"></w:document>`,
	})

	store := DocumentStore{
		Dir:      t.TempDir(),
		Route:    "/uploads",
		MaxBytes: 5 << 20,
	}

	baseName := time.Now().Format("20060102") + "_Dokumen Pengumuman.docx"
	if err := os.WriteFile(filepath.Join(store.Dir, baseName), []byte("existing"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	file := &inMemoryMultipartFile{Reader: bytes.NewReader(docxBytes)}
	fh := &multipart.FileHeader{
		Filename: "Dokumen Pengumuman.docx",
		Size:     int64(len(docxBytes)),
	}

	relativePath, err := store.SaveDocumentRelative(file, fh)
	if err != nil {
		t.Fatalf("SaveDocumentRelative returned error: %v", err)
	}

	expectedPath := "/uploads/" + time.Now().Format("20060102") + "_Dokumen Pengumuman_1.docx"
	if relativePath != expectedPath {
		t.Fatalf("expected relative path %q, got %q", expectedPath, relativePath)
	}
}

func TestSaveDocumentRelative_RejectsGenericZip(t *testing.T) {
	zipBytes := buildOOXMLZip(t, map[string]string{
		"hello.txt": "world",
	})

	store := DocumentStore{
		Dir:      t.TempDir(),
		Route:    "/uploads",
		MaxBytes: 5 << 20,
	}

	file := &inMemoryMultipartFile{Reader: bytes.NewReader(zipBytes)}
	fh := &multipart.FileHeader{
		Filename: "arsip.zip",
		Size:     int64(len(zipBytes)),
	}

	_, err := store.SaveDocumentRelative(file, fh)
	if err == nil {
		t.Fatal("expected generic zip to be rejected, got nil error")
	}

	if !strings.Contains(err.Error(), "unsupported content type: application/zip") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func buildOOXMLZip(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create zip entry %s: %v", name, err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatalf("write zip entry %s: %v", name, err)
		}
	}

	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}

	return buf.Bytes()
}
