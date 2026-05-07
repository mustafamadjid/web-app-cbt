package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	kelas_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/update"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	kelasquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	kelas_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKelasRepository struct{ mock.Mock }

func (m *MockKelasRepository) GetKelas(ctx context.Context, f kelasquery.ListKelasFilter) ([]kelas.FullKelasData, error) {
	a := m.Called(ctx, f); return a.Get(0).([]kelas.FullKelasData), a.Error(1)
}
func (m *MockKelasRepository) GetKelasById(ctx context.Context, a int, b int) (kelas.KelasData, error) {
	r := m.Called(ctx, a, b); return r.Get(0).(kelas.KelasData), r.Error(1)
}
func (m *MockKelasRepository) CreateTingkatKelas(ctx context.Context, t int) error {
	return m.Called(ctx, t).Error(0)
}
func (m *MockKelasRepository) CreateNamaKelas(ctx context.Context, n kelas.NamaKelas) error {
	return m.Called(ctx, n).Error(0)
}
func (m *MockKelasRepository) UpdateNamaKelas(ctx context.Context, id int, d updatepatch.NamaKelasPatch) error {
	return m.Called(ctx, id, d).Error(0)
}
func (m *MockKelasRepository) DeleteNamaKelas(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockKelasRepository) ExistTingkatKelas(ctx context.Context, t int) (bool, error) {
	a := m.Called(ctx, t); return a.Bool(0), a.Error(1)
}
func (m *MockKelasRepository) ExistNamaKelas(ctx context.Context, n string) (bool, error) {
	a := m.Called(ctx, n); return a.Bool(0), a.Error(1)
}

func adminCtx() context.Context {
	return middleware.WithActor(context.Background(), user.Actor{IdPengguna: 1, Role: user.ADMIN})
}

func TestUpdateNamaKelas_Success(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_update_service.NewUpdateKelasService(mockRepo)
	h := kelas_update.NewUpdateKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodPatch, "/kelas/2", strings.NewReader(`{"nama_kelas":"B"}`)).WithContext(adminCtx())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "2"}}
	mockRepo.On("UpdateNamaKelas", mock.Anything, 2, mock.MatchedBy(func(p updatepatch.NamaKelasPatch) bool {
		return p.NamaKelas != nil && *p.NamaKelas == "B"
	})).Return(nil).Once()
	h.UpdateNamaKelas(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdateNamaKelas_BadRequest_NoField(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_update_service.NewUpdateKelasService(mockRepo)
	h := kelas_update.NewUpdateKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodPatch, "/kelas/2", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "2"}}
	h.UpdateNamaKelas(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid field")
}

func TestUpdateNamaKelas_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_update_service.NewUpdateKelasService(mockRepo)
	h := kelas_update.NewUpdateKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodPatch, "/kelas/abc", strings.NewReader(`{"nama_kelas":"B"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "abc"}}
	h.UpdateNamaKelas(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id nama kelas")
}

func TestUpdateNamaKelas_NotFound(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_update_service.NewUpdateKelasService(mockRepo)
	h := kelas_update.NewUpdateKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodPatch, "/kelas/99", strings.NewReader(`{"nama_kelas":"B"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "99"}}
	mockRepo.On("UpdateNamaKelas", mock.Anything, 99, mock.Anything).Return(coreerror.ErrNotFound).Once()
	h.UpdateNamaKelas(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdateNamaKelas_BadRequest_MissingContentType(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_update_service.NewUpdateKelasService(mockRepo)
	h := kelas_update.NewUpdateKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodPatch, "/kelas/2", strings.NewReader(`{"nama_kelas":"B"}`))
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "2"}}
	h.UpdateNamaKelas(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "content type must be application/json")
}

func TestUpdateNamaKelas_BadRequest_InvalidJSON(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_update_service.NewUpdateKelasService(mockRepo)
	h := kelas_update.NewUpdateKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodPatch, "/kelas/2", strings.NewReader(`{"nama_kelas":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "2"}}
	h.UpdateNamaKelas(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}
