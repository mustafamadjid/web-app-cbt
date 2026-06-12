package httpx_test

import (
	"encoding/json"
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

func TestGetKelasIntegration(t *testing.T) {
	scope := repotestutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	admin := fixtures.CreateUser(user.ADMIN)
	testApp := handlertestutil.BuildApp(t)
	adminCookie := handlertestutil.AuthCookie(t, testApp, fixtures, admin.ID, user.ADMIN, admin.Username)

	tingkat := 900000 + int(time.Now().UnixNano()%100000)
	rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/kelas/tingkat-kelas", map[string]any{"tingkat_kelas": tingkat}, adminCookie)
	require.Equal(t, http.StatusOK, rec.Code)
	var idTingkat int64
	err := scope.Pool().QueryRow(scope.Context(), `SELECT id_kelas FROM kelas WHERE tingkat_kelas = $1`, tingkat).Scan(&idTingkat)
	require.NoError(t, err)
	scope.AddCleanupQuery(`DELETE FROM kelas WHERE id_kelas = $1`, idTingkat)

	namaKelas := "Get Integration " + repotestutil.UniqueSuffix("api")
	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodPost, "/admin/kelas/nama-kelas", map[string]any{"id_tingkat_kelas": idTingkat, "nama_kelas": namaKelas}, adminCookie)
	require.Equal(t, http.StatusOK, rec.Code)
	var idNama int64
	err = scope.Pool().QueryRow(scope.Context(), `SELECT id_nama_kelas FROM nama_kelas WHERE id_kelas = $1 AND nama_kelas = $2`, idTingkat, namaKelas).Scan(&idNama)
	require.NoError(t, err)
	scope.AddCleanupQuery(`DELETE FROM nama_kelas WHERE id_nama_kelas = $1`, idNama)

	t.Run("lists kelas successfully", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, fmt.Sprintf("/admin/kelas?q=%s&limit=10&offset=0", namaKelas), nil, adminCookie)
		assert.Equal(t, http.StatusOK, rec.Code)
		env := handlertestutil.DecodeEnvelope(t, rec)
		var data map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(env.Data, &data))
		assert.NotEmpty(t, data["item_tingkat_kelas"])
	})

	t.Run("gets kelas by ID", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, fmt.Sprintf("/admin/kelas/%d/%d", idTingkat, idNama), nil, adminCookie)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("rejects invalid query params", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, "/admin/kelas?limit=abc", nil, adminCookie)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("rejects invalid path param", func(t *testing.T) {
		rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, fmt.Sprintf("/admin/kelas/x/%d", idNama), nil, adminCookie)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
