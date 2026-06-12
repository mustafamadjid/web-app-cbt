package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonHeaderBodyValidator(t *testing.T) {
	t.Parallel()

	t.Run("invalid content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"ok":true}`))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		err := httphelper.JsonHeaderBodyValidator(rec, req)
		assert.ErrorIs(t, err, coreerror.ErrContentTypeMustJson)
	})

	t.Run("valid content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"ok":true}`))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()

		err := httphelper.JsonHeaderBodyValidator(rec, req)
		assert.NoError(t, err)
	})
}

func TestJsonHeaderBodyValidator_MaxBytesApplied(t *testing.T) {
	t.Parallel()

	tooLargeBody := strings.Repeat("a", (10<<20)+1)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tooLargeBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	err := httphelper.JsonHeaderBodyValidator(rec, req)
	require.NoError(t, err)

	_, readErr := io.ReadAll(req.Body)
	require.Error(t, readErr)
	assert.Contains(t, readErr.Error(), "http: request body too large")
}
