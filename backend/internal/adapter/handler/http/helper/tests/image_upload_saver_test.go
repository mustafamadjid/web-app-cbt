package tests

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type memoryMultipartFile struct {
	*bytes.Reader
}

func (f *memoryMultipartFile) Close() error {
	return nil
}

func TestSavePhotoRelative(t *testing.T) {
	t.Parallel()

	validPNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x60, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x01, 0xE5, 0x27, 0xD4, 0xA2, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}

	t.Run("file too large", func(t *testing.T) {
		store := httphelper.ImageStore{
			Dir:      t.TempDir(),
			Route:    "/uploads/image",
			MaxBytes: 5,
		}

		file := &memoryMultipartFile{Reader: bytes.NewReader(validPNG)}
		fh := &multipart.FileHeader{
			Filename: "avatar.png",
			Size:     int64(len(validPNG)),
		}

		_, err := store.SavePhotoRelative(file, fh)
		assert.ErrorIs(t, err, coreerror.ErrFileTooLarge)
	})

	t.Run("unsupported content type", func(t *testing.T) {
		store := httphelper.ImageStore{
			Dir:      t.TempDir(),
			Route:    "/uploads/image",
			MaxBytes: 5 << 20,
		}

		file := &memoryMultipartFile{Reader: bytes.NewReader([]byte("plain text"))}
		fh := &multipart.FileHeader{
			Filename: "avatar.txt",
			Size:     int64(len("plain text")),
		}

		_, err := store.SavePhotoRelative(file, fh)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported content type")
	})

	t.Run("success png upload", func(t *testing.T) {
		store := httphelper.ImageStore{
			Dir:      t.TempDir(),
			Route:    "/uploads/image",
			MaxBytes: 5 << 20,
		}

		file := &memoryMultipartFile{Reader: bytes.NewReader(validPNG)}
		fh := &multipart.FileHeader{
			Filename: "avatar.png",
			Size:     int64(len(validPNG)),
		}

		relativePath, err := store.SavePhotoRelative(file, fh)
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(relativePath, "/uploads/image/"))
		assert.True(t, strings.HasSuffix(relativePath, ".png"))

		savedPath := filepath.Join(store.Dir, filepath.Base(relativePath))
		_, statErr := os.Stat(savedPath)
		assert.NoError(t, statErr)
	})
}

func TestImageStore_PublicURL(t *testing.T) {
	t.Parallel()

	t.Run("missing base url", func(t *testing.T) {
		store := httphelper.ImageStore{}
		_, err := store.PublicURL("/uploads/image/a.png")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BaseURL")
	})

	t.Run("missing relative path", func(t *testing.T) {
		store := httphelper.ImageStore{BaseURL: "http://localhost:8080"}
		_, err := store.PublicURL("")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "relativePath")
	})

	t.Run("success", func(t *testing.T) {
		store := httphelper.ImageStore{BaseURL: "http://localhost:8080/"}
		url, err := store.PublicURL("/uploads/image/a.png")
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/uploads/image/a.png", url)
	})
}
