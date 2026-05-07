package httpx_test

import (
	"net/http"
	"testing"
	"time"

	handlertestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/testutil"
	repotestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateKelasIntegration(t *testing.T) {
	scope := repotestutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	admin := fixtures.CreateUser(user.ADMIN)
	siswa := fixtures.CreateUser(user.SISWA)
	testApp := handlertestutil.BuildApp(t)
	adminCookie := handlertestutil.AuthCookie(t, testApp, fixtures, admin.ID, user.ADMIN, admin.Username)
	siswaCookie := handlertestutil.AuthCookie(t, testApp, fixtures, siswa.ID, user.SISWA, siswa.Username)

	t.Run("rejects unauthorized request", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/kelas/tingkat-kelas", map[string]any{"tingkat_kelas": 1}, nil)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("rejects forbidden role", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/kelas/tingkat-kelas", map[string]any{"tingkat_kelas": 1}, siswaCookie)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("validates create tingkat kelas", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/kelas/tingkat-kelas", map[string]any{"tingkat_kelas": 0}, adminCookie)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		env := handlertestutil.DecodeEnvelope(t, rec)
		require.NotNil(t, env.Error)
		assert.Equal(t, "BAD_REQUEST", env.Error.Code)
	})

	t.Run("creates tingkat kelas and nama kelas successfully", func(t *testing.T) {
		tingkat := 900000 + int(time.Now().UnixNano()%100000)
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/kelas/tingkat-kelas", map[string]any{"tingkat_kelas": tingkat}, adminCookie)
		assert.Equal(t, http.StatusOK, rec.Code)

		var idTingkat int64
		err := scope.Pool().QueryRow(scope.Context(), `SELECT id_kelas FROM kelas WHERE tingkat_kelas = $1`, tingkat).Scan(&idTingkat)
		require.NoError(t, err)
		scope.AddCleanupQuery(`DELETE FROM kelas WHERE id_kelas = $1`, idTingkat)

		namaKelas := "Create Integration " + repotestutil.UniqueSuffix("api")
		rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/kelas/nama-kelas", map[string]any{
			"id_tingkat_kelas": idTingkat,
			"nama_kelas":       "  " + namaKelas + "  ",
		}, adminCookie)
		assert.Equal(t, http.StatusOK, rec.Code)

		var idNama int64
		err = scope.Pool().QueryRow(scope.Context(), `SELECT id_nama_kelas FROM nama_kelas WHERE id_kelas = $1 AND nama_kelas = $2`, idTingkat, namaKelas).Scan(&idNama)
		require.NoError(t, err)
		scope.AddCleanupQuery(`DELETE FROM nama_kelas WHERE id_nama_kelas = $1`, idNama)
	})
}
