package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCORSPolicy_AllowsTrustedOrigin(t *testing.T) {
	t.Setenv("TRUSTED_ORIGINS", "https://app.example.com, https://admin.example.com")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodOptions, "https://app.example.com/api", nil)
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec := httptest.NewRecorder()

	middleware.CORSPolicy(next).ServeHTTP(rec, req)

	assert.Equal(t, "https://app.example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestPreventCSRF_AllowsTrustedOrigin(t *testing.T) {
	t.Setenv("TRUSTED_ORIGINS", "https://app.example.com")

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "https://app.example.com/api", nil)
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Referer", "https://app.example.com/form")
	rec := httptest.NewRecorder()

	middleware.PreventCSRF(next).ServeHTTP(rec, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestPreventCSRF_DeniesUntrustedOrigin(t *testing.T) {
	t.Setenv("TRUSTED_ORIGINS", "https://app.example.com")

	logger := &capturedLogger{}
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodPost, "https://app.example.com/api", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	req.Header.Set("Referer", "https://evil.example.com/form")
	req = req.WithContext(corelog.WithLogger(req.Context(), logger))
	rec := httptest.NewRecorder()

	middleware.PreventCSRF(next).ServeHTTP(rec, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	require.Len(t, logger.errorCalls, 1)
	assert.Equal(t, "csrf origin denied", logger.errorCalls[0].msg)
}
