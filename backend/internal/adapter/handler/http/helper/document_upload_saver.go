package httpx

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

type DocumentStore struct {
	Dir      string // contoh: "./public/uploads"
	BaseURL  string // contoh: "https://example.com"
	Route    string // contoh: "/uploads" (route publik)
	MaxBytes int64  // contoh: 5 << 20
}

func (s *DocumentStore) SaveDocumentRelative(file multipart.File, fh *multipart.FileHeader) (string, error) {
	if s.Dir == "" || s.Route == "" {
		return "", errors.New("DocumentStore Dir/Route must be set")
	}

	if s.MaxBytes <= 0 {
		s.MaxBytes = 5 << 20
	}

	if fh != nil && fh.Size > s.MaxBytes {
		return "", coreerror.ErrFileTooLarge
	}

	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		var pe *os.PathError
		if errors.As(err, &pe) {
			switch {
			case errors.Is(pe.Err, syscall.EACCES), errors.Is(pe.Err, syscall.EPERM):
				return "", errors.New("mkdir: permission denied")
			case errors.Is(pe.Err, syscall.EROFS):
				return "", errors.New("mkdir: read-only file system")
			}
		}
		return "", err
	}
	const sniffLen = 512
	head := make([]byte, sniffLen)

	n, err := io.ReadFull(file, head)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", fmt.Errorf("read head: %w", err)
	}
	head = head[:n]

	mime := detectDocumentMIME(file, fh, head)
	ext, ok := allowedDocumentExt(mime)
	if !ok {
		return "", fmt.Errorf("unsupported content type: %s", mime)
	}

	seekable := false
	if seeker, ok := file.(io.Seeker); ok {
		seekable = true
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return "", fmt.Errorf("seek: %w", err)
		}
	}

	name := randomHex(16)
	datePrefix := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s%s", datePrefix, name, ext)

	relativePath := strings.TrimRight(s.Route, "/") + "/" + filename

	dstPath := filepath.Join(s.Dir, filename)

	tmpPath := dstPath + ".tmp"
	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return "", fmt.Errorf("open tmp: %w", err)
	}

	success := false
	defer func() {
		_ = out.Close()
		if !success {
			_ = os.Remove(tmpPath)
		}
	}()

	var written int64

	if seekable {
		lr := &io.LimitedReader{R: file, N: s.MaxBytes + 1}
		nn, err := io.Copy(out, lr)
		written = nn
		if err != nil {
			return "", fmt.Errorf("copy: %w", err)
		}
		if lr.N <= 0 {
			return "", fmt.Errorf("file exceeds max bytes: %d", s.MaxBytes)
		}
	} else {
		if len(head) == 0 {
			return "", errors.New("empty file")
		}
		nn, err := out.Write(head)
		if err != nil {
			return "", fmt.Errorf("write head: %w", err)
		}
		written += int64(nn)

		lr := &io.LimitedReader{R: file, N: (s.MaxBytes - written) + 1}
		nn2, err := io.Copy(out, lr)
		written += nn2
		if err != nil {
			return "", fmt.Errorf("copy rest: %w", err)
		}
		if lr.N <= 0 {
			return "", fmt.Errorf("file exceeds max bytes: %d", s.MaxBytes)
		}
	}

	if written == 0 {
		return "", errors.New("empty file")
	}

	if err := out.Close(); err != nil {
		return "", fmt.Errorf("close: %w", err)
	}

	if err := os.Rename(tmpPath, dstPath); err != nil {
		return "", fmt.Errorf("rename: %w", err)
	}

	success = true
	return relativePath, nil
}

func allowedDocumentExt(mime string) (string, bool) {
	if mime == "application/pdf" {
		return ".pdf", true
	}

	if mime == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		return ".docx", true
	}
	return "", false
}

func detectDocumentMIME(file multipart.File, fh *multipart.FileHeader, head []byte) string {
	mime := http.DetectContentType(head)
	if mime != "application/zip" {
		return mime
	}

	if isDocxOOXML(file, fh) {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}

	return mime
}

func isDocxOOXML(file multipart.File, fh *multipart.FileHeader) bool {
	size, ok := resolveUploadFileSize(file, fh)
	if !ok || size <= 0 {
		return false
	}

	zr, err := zip.NewReader(file, size)
	if err != nil {
		return false
	}

	hasContentTypes := false
	hasWordDocument := false

	for _, f := range zr.File {
		switch f.Name {
		case "[Content_Types].xml":
			hasContentTypes = true
		case "word/document.xml":
			hasWordDocument = true
		}

		if hasContentTypes && hasWordDocument {
			return true
		}
	}

	return false
}

func resolveUploadFileSize(file multipart.File, fh *multipart.FileHeader) (int64, bool) {
	if fh != nil && fh.Size > 0 {
		return fh.Size, true
	}

	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, false
	}

	endPos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		_, _ = file.Seek(currentPos, io.SeekStart)
		return 0, false
	}

	if _, err := file.Seek(currentPos, io.SeekStart); err != nil {
		return 0, false
	}

	return endPos, true
}
