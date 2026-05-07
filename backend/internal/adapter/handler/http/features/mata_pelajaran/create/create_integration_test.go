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

func TestCreateMapelIntegration(t *testing.T) {
	scope := repotestutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	admin := fixtures.CreateUser(user.ADMIN)
	siswa := fixtures.CreateUser(user.SISWA)
	kelas := fixtures.CreateKelas(10)
	testApp := handlertestutil.BuildApp(t)
	adminCookie := handlertestutil.AuthCookie(t, testApp, fixtures, admin.ID, user.ADMIN, admin.Username)
	siswaCookie := handlertestutil.AuthCookie(t, testApp, fixtures, siswa.ID, user.SISWA, siswa.Username)

	t.Run("rejects unauthorized", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/mata-pelajaran", map[string]any{"id_kelas": kelas.ID}, nil)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
	t.Run("rejects forbidden role", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/mata-pelajaran", map[string]any{"id_kelas": kelas.ID}, siswaCookie)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
	t.Run("validates create request", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/mata-pelajaran", map[string]any{"id_kelas": kelas.ID}, adminCookie)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("creates mapel successfully", func(t *testing.T) {
		kode := fmt.Sprintf("MPAPI%d", time.Now().UnixNano()%1_000_000_000)
		nama := "Mapel Create " + repotestutil.UniqueSuffix("api")
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/mata-pelajaran", map[string]any{
			"id_kelas": kelas.ID, "kode_mapel": kode, "nama_mapel": nama, "deskripsi": "Deskripsi",
		}, adminCookie)
		assert.Equal(t, http.StatusOK, rec.Code)
		var idMapel int64
		err := scope.Pool().QueryRow(scope.Context(), `SELECT id_mapel FROM mata_pelajaran WHERE kode_mapel = $1`, kode).Scan(&idMapel)
		require.NoError(t, err)
		scope.AddCleanupQuery(`DELETE FROM mata_pelajaran WHERE id_mapel = $1`, idMapel)
	})
}
