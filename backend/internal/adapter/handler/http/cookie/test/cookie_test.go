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

func TestAuthCookieOperations(t *testing.T) {
	exp := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	cfg := cookie.CookieConfig{AccessName: "access_token", RefreshName: "refresh_token", Secure: true, SameSite: http.SameSiteStrictMode}
	tests := []struct {
		name, operation string
		wantCount       int
		wantNames       []string
		wantPaths       []string
		wantValues      []string
		cleared         bool
	}{
		{name: "Path 1 -> access cookie uses root path and access configuration", operation: "access", wantCount: 1, wantNames: []string{"access_token"}, wantPaths: []string{"/"}, wantValues: []string{"access-value"}},
		{name: "Path 2 -> refresh cookie is scoped to auth path", operation: "refresh", wantCount: 1, wantNames: []string{"refresh_token"}, wantPaths: []string{"/auth"}, wantValues: []string{"refresh-value"}},
		{name: "Path 3 -> clear emits expired access and refresh cookies", operation: "clear", wantCount: 2, wantNames: []string{"access_token", "refresh_token"}, wantPaths: []string{"/", "/auth"}, wantValues: []string{"", ""}, cleared: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			switch tt.operation {
			case "access":
				cookie.SetAccessCookie(rec, cfg, tt.wantValues[0], exp)
			case "refresh":
				cookie.SetRefreshCookie(rec, cfg, tt.wantValues[0], exp)
			case "clear":
				cookie.ClearAuthCookies(rec, cfg)
			}
			got := rec.Result().Cookies()
			require.Len(t, got, tt.wantCount)
			for i := range got {
				assert.Equal(t, tt.wantNames[i], got[i].Name)
				assert.Equal(t, tt.wantPaths[i], got[i].Path)
				assert.Equal(t, tt.wantValues[i], got[i].Value)
				assert.True(t, got[i].HttpOnly)
				assert.True(t, got[i].Secure)
				assert.Equal(t, http.SameSiteStrictMode, got[i].SameSite)
				if tt.cleared {
					assert.Equal(t, -1, got[i].MaxAge)
					assert.Equal(t, int64(0), got[i].Expires.Unix())
				} else {
					assert.Equal(t, exp, got[i].Expires)
				}
			}
		})
	}
}
