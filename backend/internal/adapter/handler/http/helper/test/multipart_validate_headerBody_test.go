package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultipartHeaderBodyValidator(t *testing.T) {
	t.Parallel()

	t.Run("invalid content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("plain"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		err := httphelper.MultipartHeaderBodyValidator(rec, req, 10<<20)
		assert.ErrorIs(t, err, coreerror.ErrContentTypeMustMultipart)
	})

	t.Run("invalid multipart form", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("broken body"))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=bad")
		rec := httptest.NewRecorder()

		err := httphelper.MultipartHeaderBodyValidator(rec, req, 10<<20)
		assert.ErrorIs(t, err, coreerror.ErrInvalidMultipartForm)
	})

	t.Run("valid multipart form", func(t *testing.T) {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		require.NoError(t, writer.WriteField("judul", "pengumuman"))
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		err := httphelper.MultipartHeaderBodyValidator(rec, req, 10<<20)
		require.NoError(t, err)
		assert.Equal(t, "pengumuman", req.FormValue("judul"))
	})

	t.Run("uses default max bytes when maxBytes <= 0", func(t *testing.T) {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		require.NoError(t, writer.WriteField("field", "value"))
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		err := httphelper.MultipartHeaderBodyValidator(rec, req, 0)
		require.NoError(t, err)
		assert.Equal(t, "value", req.FormValue("field"))
	})
}
