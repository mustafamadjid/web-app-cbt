package cookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cookie "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAccessCookie(t *testing.T) {
	exp := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	rec := httptest.NewRecorder()

	cookie.SetAccessCookie(rec, cookie.CookieConfig{
		AccessName: "access_token",
		Secure:     true,
		SameSite:   http.SameSiteStrictMode,
	}, "access-value", exp)

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 1)
	got := cookies[0]
	assert.Equal(t, "access_token", got.Name)
	assert.Equal(t, "access-value", got.Value)
	assert.Equal(t, "/", got.Path)
	assert.True(t, got.HttpOnly)
	assert.True(t, got.Secure)
	assert.Equal(t, http.SameSiteStrictMode, got.SameSite)
	assert.Equal(t, exp, got.Expires)
}

func TestSetRefreshCookie(t *testing.T) {
	exp := time.Date(2026, time.February, 3, 4, 5, 6, 0, time.UTC)
	rec := httptest.NewRecorder()

	cookie.SetRefreshCookie(rec, cookie.CookieConfig{
		RefreshName: "refresh_token",
		Secure:      true,
		SameSite:    http.SameSiteLaxMode,
	}, "refresh-value", exp)

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 1)
	got := cookies[0]
	assert.Equal(t, "refresh_token", got.Name)
	assert.Equal(t, "refresh-value", got.Value)
	assert.Equal(t, "/auth", got.Path)
	assert.True(t, got.HttpOnly)
	assert.True(t, got.Secure)
	assert.Equal(t, http.SameSiteLaxMode, got.SameSite)
	assert.Equal(t, exp, got.Expires)
}

func TestClearAuthCookies(t *testing.T) {
	rec := httptest.NewRecorder()

	cookie.ClearAuthCookies(rec, cookie.CookieConfig{
		AccessName:  "access_token",
		RefreshName: "refresh_token",
		Secure:      true,
		SameSite:    http.SameSiteNoneMode,
	})

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 2)
	assert.Equal(t, "access_token", cookies[0].Name)
	assert.Equal(t, "", cookies[0].Value)
	assert.Equal(t, "/", cookies[0].Path)
	assert.Equal(t, -1, cookies[0].MaxAge)
	assert.Equal(t, int64(0), cookies[0].Expires.Unix())
	assert.True(t, cookies[0].HttpOnly)
	assert.True(t, cookies[0].Secure)
	assert.Equal(t, http.SameSiteNoneMode, cookies[0].SameSite)

	assert.Equal(t, "refresh_token", cookies[1].Name)
	assert.Equal(t, "", cookies[1].Value)
	assert.Equal(t, "/auth", cookies[1].Path)
	assert.Equal(t, -1, cookies[1].MaxAge)
	assert.Equal(t, int64(0), cookies[1].Expires.Unix())
	assert.True(t, cookies[1].HttpOnly)
	assert.True(t, cookies[1].Secure)
	assert.Equal(t, http.SameSiteNoneMode, cookies[1].SameSite)
}
