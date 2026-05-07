package integration_test

import (
	"net/http"
	"testing"

	handlertestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/testutil"
	repotestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardHandlerIntegration(t *testing.T) {
	scope := repotestutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	admin := fixtures.CreateUser(user.ADMIN)
	guru := fixtures.CreateUser(user.GURU)
	siswa := fixtures.CreateUser(user.SISWA)
	testApp := handlertestutil.BuildApp(t)

	adminCookie := handlertestutil.AuthCookie(t, testApp, fixtures, admin.ID, user.ADMIN, admin.Username)
	guruCookie := handlertestutil.AuthCookie(t, testApp, fixtures, guru.ID, user.GURU, guru.Username)
	siswaCookie := handlertestutil.AuthCookie(t, testApp, fixtures, siswa.ID, user.SISWA, siswa.Username)

	rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, "/admin/dashboard", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, "/admin/dashboard", nil, siswaCookie)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, "/admin/dashboard", nil, adminCookie)
	assert.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, handlertestutil.DecodeEnvelope(t, rec).Error)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, "/guru/dashboard", nil, guruCookie)
	assert.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, handlertestutil.DecodeEnvelope(t, rec).Error)
}
