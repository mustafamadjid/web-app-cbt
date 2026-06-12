package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	pengumuman_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/get"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPengumumanRepo struct{ mock.Mock }

func (m *MockPengumumanRepo) GetPengumumanActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx); return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *MockPengumumanRepo) GetPengumumanNonActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx); return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *MockPengumumanRepo) GetPengumumanIncoming(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	a := m.Called(ctx); return a.Get(0).([]pengumuman.Pengumuman), a.Error(1)
}
func (m *MockPengumumanRepo) GetPengumumanById(ctx context.Context, id pengumuman.ID) (pengumuman.Pengumuman, error) {
	a := m.Called(ctx, id); return a.Get(0).(pengumuman.Pengumuman), a.Error(1)
}
func (m *MockPengumumanRepo) CreatePengumuman(ctx context.Context, p pengumuman.Pengumuman) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockPengumumanRepo) UpdatePengumuman(ctx context.Context, id pengumuman.ID, p updatepatch.PengumumanUpdatePatch) error {
	return m.Called(ctx, id, p).Error(0)
}
func (m *MockPengumumanRepo) DeletePengumuman(ctx context.Context, id pengumuman.ID) error {
	return m.Called(ctx, id).Error(0)
}

func TestGetPengumumanActive_Success(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	svc := pengumuman_get_service.NewGetPengumumanService(mockRepo)
	h := pengumuman_get.NewGetPengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengumuman/active", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetPengumumanActive", mock.Anything).Return([]pengumuman.Pengumuman{
		{IdPengumuman: 1, JudulPengumuman: "Test"},
	}, nil).Once()
	h.GetPengumumanActive(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test")
}

func TestGetPengumumanNonActive_Success(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	svc := pengumuman_get_service.NewGetPengumumanService(mockRepo)
	h := pengumuman_get.NewGetPengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengumuman/non-active", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetPengumumanNonActive", mock.Anything).Return([]pengumuman.Pengumuman{}, nil).Once()
	h.GetPengumumanNonActive(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPengumumanIncoming_Success(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	svc := pengumuman_get_service.NewGetPengumumanService(mockRepo)
	h := pengumuman_get.NewGetPengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengumuman/incoming", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetPengumumanIncoming", mock.Anything).Return([]pengumuman.Pengumuman{}, nil).Once()
	h.GetPengumumanIncoming(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPengumumanByID_Success(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	svc := pengumuman_get_service.NewGetPengumumanService(mockRepo)
	h := pengumuman_get.NewGetPengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengumuman/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengumuman", Value: "1"}}
	mockRepo.On("GetPengumumanById", mock.Anything, pengumuman.ID(1)).Return(pengumuman.Pengumuman{IdPengumuman: 1, JudulPengumuman: "Test"}, nil).Once()
	h.GetPengumumanByID(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test")
}

func TestGetPengumumanByID_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	svc := pengumuman_get_service.NewGetPengumumanService(mockRepo)
	h := pengumuman_get.NewGetPengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengumuman/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengumuman", Value: "abc"}}
	h.GetPengumumanByID(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id pengumuman")
}

func TestGetPengumumanByID_NotFound(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	svc := pengumuman_get_service.NewGetPengumumanService(mockRepo)
	h := pengumuman_get.NewGetPengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengumuman/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengumuman", Value: "99"}}
	mockRepo.On("GetPengumumanById", mock.Anything, pengumuman.ID(99)).Return(pengumuman.Pengumuman{}, coreerror.ErrNotFound).Once()
	h.GetPengumumanByID(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPengumumanActive_InternalServerError(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	svc := pengumuman_get_service.NewGetPengumumanService(mockRepo)
	h := pengumuman_get.NewGetPengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengumuman/active", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetPengumumanActive", mock.Anything).Return([]pengumuman.Pengumuman{}, coreerror.ErrDbError).Once()
	h.GetPengumumanActive(w, req, nil)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
