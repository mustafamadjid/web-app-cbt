package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	middleware "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardRateLimit_BypassesWhenNoActor(t *testing.T) {
	limiter := &fakeRateLimiter{allowed: true}
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.StandardRateLimit(limiter, next).ServeHTTP(rec, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 0, limiter.calls)
}

func TestStandardRateLimit_BlockedRequestSetsRetryAfter(t *testing.T) {
	limiter := &fakeRateLimiter{
		allowed:    false,
		retryAfter: 1500 * time.Millisecond,
	}

	nextCalled := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{
		IdPengguna: 3,
		Role:       user.ADMIN,
		Username:   "Alice",
	}))
	rec := httptest.NewRecorder()

	middleware.StandardRateLimit(limiter, next).ServeHTTP(rec, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Equal(t, "username:Alice:203.0.113.10", limiter.key)
	assert.Equal(t, "2", rec.Header().Get("Retry-After"))
}

func TestLoginRateLimit_AllowsRequestAndRestoresBody(t *testing.T) {
	limiter := &fakeRateLimiter{allowed: true}
	originalBody := []byte(`{"username":" Alice ","password":"secret"}`)

	var bodySeen []byte
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		bodySeen, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(originalBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "198.51.100.7")
	rec := httptest.NewRecorder()

	middleware.LoginRateLimit(limiter, next).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "login:alice:198.51.100.7", limiter.key)
	assert.Equal(t, originalBody, bodySeen)
}

func TestLoginRateLimit_InvalidJSONFallsBackToAnonymousKey(t *testing.T) {
	limiter := &fakeRateLimiter{allowed: true}
	nextCalled := false
	invalidBody := []byte(`{"username"`)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		bodySeen, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, invalidBody, bodySeen)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "198.51.100.7")
	rec := httptest.NewRecorder()

	middleware.LoginRateLimit(limiter, next).ServeHTTP(rec, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "login:_nouser_:198.51.100.7", limiter.key)
	assert.Equal(t, 1, limiter.calls)
}
