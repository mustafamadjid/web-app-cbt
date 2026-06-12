package tests

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreFileToDisk(t *testing.T) {
	t.Parallel()

	t.Run("required file missing", func(t *testing.T) {
		req := newMultipartRequestWithoutFile(t)

		path, err := httphelper.StoreFileToDisk(req, "foto", true, func(file multipart.File, fh *multipart.FileHeader) (string, error) {
			return "", nil
		})

		assert.Nil(t, path)
		assert.ErrorIs(t, err, coreerror.ErrMissingField)
	})

	t.Run("optional file missing", func(t *testing.T) {
		req := newMultipartRequestWithoutFile(t)
		called := false

		path, err := httphelper.StoreFileToDisk(req, "foto", false, func(file multipart.File, fh *multipart.FileHeader) (string, error) {
			called = true
			return "", nil
		})

		assert.NoError(t, err)
		assert.Nil(t, path)
		assert.False(t, called)
	})

	t.Run("invalid multipart body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid multipart"))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=broken")

		path, err := httphelper.StoreFileToDisk(req, "foto", false, func(file multipart.File, fh *multipart.FileHeader) (string, error) {
			return "", nil
		})

		assert.Nil(t, path)
		assert.Error(t, err)
		assert.False(t, errors.Is(err, http.ErrMissingFile))
	})

	t.Run("saver returns error", func(t *testing.T) {
		req := newMultipartRequestWithFile(t, "foto", "avatar.png", []byte("abc"))
		saveErr := errors.New("save failed")

		path, err := httphelper.StoreFileToDisk(req, "foto", false, func(file multipart.File, fh *multipart.FileHeader) (string, error) {
			return "", saveErr
		})

		assert.Nil(t, path)
		assert.ErrorIs(t, err, saveErr)
	})

	t.Run("success", func(t *testing.T) {
		req := newMultipartRequestWithFile(t, "foto", "avatar.png", []byte("abc"))
		called := false

		path, err := httphelper.StoreFileToDisk(req, "foto", false, func(file multipart.File, fh *multipart.FileHeader) (string, error) {
			called = true
			require.Equal(t, "avatar.png", fh.Filename)
			_, readErr := io.ReadAll(file)
			require.NoError(t, readErr)
			return "/uploads/image/avatar.png", nil
		})

		require.NoError(t, err)
		require.NotNil(t, path)
		assert.True(t, called)
		assert.Equal(t, "/uploads/image/avatar.png", *path)
	})
}

func newMultipartRequestWithoutFile(t *testing.T) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("name", "tester"))
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func newMultipartRequestWithFile(t *testing.T, fieldName, fileName string, content []byte) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
