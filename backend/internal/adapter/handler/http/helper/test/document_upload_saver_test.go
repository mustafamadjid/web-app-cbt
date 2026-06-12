package tests

import (
	"archive/zip"
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
	"time"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type inMemoryMultipartFile struct {
	*bytes.Reader
}

func (f *inMemoryMultipartFile) Close() error {
	return nil
}

func TestSaveDocumentRelative_AcceptsDocxOOXMLZip(t *testing.T) {
	t.Parallel()

	docxBytes := buildOOXMLZip(t, map[string]string{
		"[Content_Types].xml": `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"></Types>`,
		"word/document.xml":   `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"></w:document>`,
	})

	store := httphelper.DocumentStore{
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
	require.NoError(t, err)

	expectedFileName := time.Now().Format("20060102") + "_Surat Edaran UTS.docx"
	assert.Equal(t, "/uploads/"+expectedFileName, relativePath)

	savedPath := filepath.Join(store.Dir, filepath.Base(relativePath))
	_, statErr := os.Stat(savedPath)
	assert.NoError(t, statErr)
}

func TestSaveDocumentRelative_UsesIncrementSuffixWhenFileExists(t *testing.T) {
	t.Parallel()

	docxBytes := buildOOXMLZip(t, map[string]string{
		"[Content_Types].xml": `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"></Types>`,
		"word/document.xml":   `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"></w:document>`,
	})

	store := httphelper.DocumentStore{
		Dir:      t.TempDir(),
		Route:    "/uploads",
		MaxBytes: 5 << 20,
	}

	baseName := time.Now().Format("20060102") + "_Dokumen Pengumuman.docx"
	require.NoError(t, os.WriteFile(filepath.Join(store.Dir, baseName), []byte("existing"), 0o644))

	file := &inMemoryMultipartFile{Reader: bytes.NewReader(docxBytes)}
	fh := &multipart.FileHeader{
		Filename: "Dokumen Pengumuman.docx",
		Size:     int64(len(docxBytes)),
	}

	relativePath, err := store.SaveDocumentRelative(file, fh)
	require.NoError(t, err)

	expectedPath := "/uploads/" + time.Now().Format("20060102") + "_Dokumen Pengumuman_1.docx"
	assert.Equal(t, expectedPath, relativePath)
}

func TestSaveDocumentRelative_RejectsGenericZip(t *testing.T) {
	t.Parallel()

	zipBytes := buildOOXMLZip(t, map[string]string{
		"hello.txt": "world",
	})

	store := httphelper.DocumentStore{
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
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported content type: application/zip")
}

func buildOOXMLZip(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for name, body := range files {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(body))
		require.NoError(t, err)
	}

	require.NoError(t, zw.Close())
	return buf.Bytes()
}
