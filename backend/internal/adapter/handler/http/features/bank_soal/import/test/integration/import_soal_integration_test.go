package integration_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	handlertestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/testutil"
	repotestutil "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	importdomain "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportSoalHandlerIntegration(t *testing.T) {
	scope := repotestutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	guru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guru.ID)
	kelas := fixtures.CreateKelas(10)
	mapel := fixtures.CreateMapel(kelas.ID)
	bank := fixtures.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)
	job := fixtures.CreateImportJob(bank.ID, guru.ID, importdomain.StatusPending, "fixtures/import.docx")
	testApp := handlertestutil.BuildApp(t)
	cookie := handlertestutil.AuthCookie(t, testApp, fixtures, guru.ID, user.GURU, guru.Username)

	rec := handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, fmt.Sprintf("/admin/bank-soal/import-job/%d", job.ID), nil, nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	rec = handlertestutil.DoRaw(t, testApp.HTTP.Handler, http.MethodPost, fmt.Sprintf("/admin/bank-soal/import/%d", bank.ID), strings.NewReader(`{}`), "application/json", cookie)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, "/admin/bank-soal/import-job/abc", nil, cookie)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, fmt.Sprintf("/admin/bank-soal/import-job/%d", job.ID), nil, cookie)
	assert.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, handlertestutil.DecodeEnvelope(t, rec).Error)

	rec = handlertestutil.DoJSON(t, testApp.HTTP.Handler, http.MethodGet, fmt.Sprintf("/admin/bank-soal/import-jobs/%d", bank.ID), nil, cookie)
	assert.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, handlertestutil.DecodeEnvelope(t, rec).Error)
}
