package integration_test

import (
	"net/http"
	"testing"
	"time"

	handlertestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/testutil"
	repotestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/security/bcrypt"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandlerIntegration_APIFlow(t *testing.T) {
	scope := repotestutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	loginUser := fixtures.CreateUser(user.GURU)
	admin := fixtures.CreateUser(user.ADMIN)
	testApp := handlertestutil.BuildApp(t)

	hash, err := bcrypt.NewHasher(4).GenerateHash("password-login")
	require.NoError(t, err)
	_, err = scope.Pool().Exec(scope.Context(), `UPDATE pengguna SET password = $1 WHERE id_pengguna = $2`, hash, loginUser.ID)
	require.NoError(t, err)

	rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/auth/login", map[string]any{
		"username": loginUser.Username,
	}, nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/auth/login", map[string]any{
		"username": loginUser.Username,
		"password": "wrong-password",
	}, nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/auth/login", map[string]any{
		"username": loginUser.Username,
		"password": "password-login",
	}, nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, handlertestutil.DecodeEnvelope(t, rec).Error)

	accessCookie := findCookie(t, rec.Result().Cookies(), "access_token")
	refreshCookie := findCookie(t, rec.Result().Cookies(), "refresh_token")

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, "/auth/me", nil, accessCookie)
	assert.Equal(t, http.StatusOK, rec.Code)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/auth/refresh", nil, refreshCookie)
	assert.Equal(t, http.StatusOK, rec.Code)

	req := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/auth/logout", nil, refreshCookie)
	req.Result().Cookies()
	assert.Equal(t, http.StatusOK, req.Code)

	adminCookie := handlertestutil.AuthCookie(t, testApp, fixtures, admin.ID, user.ADMIN, admin.Username)
	session := fixtures.CreateSession(loginUser.ID, user.GURU, time.Now().Add(time.Hour), nil)
	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPut, "/admin/auth/revoke-session", map[string]any{
		"session_id": session.ID,
	}, adminCookie)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func findCookie(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("cookie %s not found", name)
	return nil
}
