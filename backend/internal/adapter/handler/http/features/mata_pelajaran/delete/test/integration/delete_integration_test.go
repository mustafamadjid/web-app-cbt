package httpx_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	handlertestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/testutil"
	repotestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteMapelIntegration(t *testing.T) {
	scope := repotestutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	admin := fixtures.CreateUser(user.ADMIN)
	kelas := fixtures.CreateKelas(10)
	testApp := handlertestutil.BuildApp(t)
	adminCookie := handlertestutil.AuthCookie(t, testApp, fixtures, admin.ID, user.ADMIN, admin.Username)

	kode := fmt.Sprintf("MPDEL%d", time.Now().UnixNano()%1_000_000_000)
	nama := "Mapel Delete " + repotestutil.UniqueSuffix("api")
	rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/mata-pelajaran", map[string]any{
		"id_kelas": kelas.ID, "kode_mapel": kode, "nama_mapel": nama, "deskripsi": "Deskripsi",
	}, adminCookie)
	require.Equal(t, http.StatusOK, rec.Code)
	var idMapel int64
	err := scope.Pool().QueryRow(scope.Context(), `SELECT id_mapel FROM mata_pelajaran WHERE kode_mapel = $1`, kode).Scan(&idMapel)
	require.NoError(t, err)

	t.Run("deletes mapel and returns not found", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodDelete, fmt.Sprintf("/admin/mata-pelajaran/%d", idMapel), nil, adminCookie)
		assert.Equal(t, http.StatusOK, rec.Code)
		rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, fmt.Sprintf("/admin/mata-pelajaran/%d", idMapel), nil, adminCookie)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
