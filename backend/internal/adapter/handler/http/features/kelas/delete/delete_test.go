package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	kelas_delete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/delete"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	kelasquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	kelas_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/delete"
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

func TestDeleteKelas_Success(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_delete_service.NewDeleteKelasService(mockRepo)
	h := kelas_delete.NewDeleteKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodDelete, "/kelas/2", nil).WithContext(adminCtx())
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "2"}}
	mockRepo.On("DeleteNamaKelas", mock.Anything, 2).Return(nil).Once()
	h.DeleteKelas(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeleteKelas_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_delete_service.NewDeleteKelasService(mockRepo)
	h := kelas_delete.NewDeleteKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodDelete, "/kelas/0", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "0"}}
	h.DeleteKelas(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id nama kelas")
}

func TestDeleteKelas_BadRequest_NonNumericParam(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_delete_service.NewDeleteKelasService(mockRepo)
	h := kelas_delete.NewDeleteKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodDelete, "/kelas/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "abc"}}
	h.DeleteKelas(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id nama kelas")
}

func TestDeleteKelas_NotFound(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_delete_service.NewDeleteKelasService(mockRepo)
	h := kelas_delete.NewDeleteKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodDelete, "/kelas/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "99"}}
	mockRepo.On("DeleteNamaKelas", mock.Anything, 99).Return(coreerror.ErrNotFound).Once()
	h.DeleteKelas(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeleteKelas_Restricted(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_delete_service.NewDeleteKelasService(mockRepo)
	h := kelas_delete.NewDeleteKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodDelete, "/kelas/3", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "3"}}
	mockRepo.On("DeleteNamaKelas", mock.Anything, 3).Return(coreerror.ErrDeleteRestricted).Once()
	h.DeleteKelas(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "DELETE_RESTRICTED")
	mockRepo.AssertExpectations(t)
}

func TestDeleteKelas_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	svc := kelas_delete_service.NewDeleteKelasService(mockRepo)
	h := kelas_delete.NewDeleteKelasHandler(svc, nil)
	req := httptest.NewRequest(http.MethodGet, "/kelas/2", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idNamaKelas", Value: "2"}}
	h.DeleteKelas(w, req, ps)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
