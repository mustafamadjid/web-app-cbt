package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	pengumuman_delete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/delete"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/delete"
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

type MockDeleteFileRepo struct{ mock.Mock }

func (m *MockDeleteFileRepo) DeleteFile(ctx context.Context, filePath string) error {
	return m.Called(ctx, filePath).Error(0)
}

func TestDeletePengumuman_Success(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	mockDeleteFile := new(MockDeleteFileRepo)
	svc := pengumuman_delete_service.NewDeletePengumumanService(mockRepo, mockDeleteFile)
	h := pengumuman_delete.NewDeletePengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/pengumuman/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengumuman", Value: "1"}}
	mockRepo.On("GetPengumumanById", mock.Anything, pengumuman.ID(1)).Return(pengumuman.Pengumuman{IdPengumuman: 1, DokumenPengumuman: ""}, nil).Once()
	mockDeleteFile.On("DeleteFile", mock.Anything, mock.Anything).Return(nil).Maybe()
	mockRepo.On("DeletePengumuman", mock.Anything, pengumuman.ID(1)).Return(nil).Once()
	h.DeletePengumuman(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeletePengumuman_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	mockDeleteFile := new(MockDeleteFileRepo)
	svc := pengumuman_delete_service.NewDeletePengumumanService(mockRepo, mockDeleteFile)
	h := pengumuman_delete.NewDeletePengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/pengumuman/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengumuman", Value: "abc"}}
	h.DeletePengumuman(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id pengumuman")
}

func TestDeletePengumuman_NotFound(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	mockDeleteFile := new(MockDeleteFileRepo)
	svc := pengumuman_delete_service.NewDeletePengumumanService(mockRepo, mockDeleteFile)
	h := pengumuman_delete.NewDeletePengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/pengumuman/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengumuman", Value: "99"}}
	mockRepo.On("GetPengumumanById", mock.Anything, pengumuman.ID(99)).Return(pengumuman.Pengumuman{}, coreerror.ErrNotFound).Once()
	mockRepo.On("DeletePengumuman", mock.Anything, pengumuman.ID(99)).Return(coreerror.ErrNotFound).Once()
	h.DeletePengumuman(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeletePengumuman_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	mockDeleteFile := new(MockDeleteFileRepo)
	svc := pengumuman_delete_service.NewDeletePengumumanService(mockRepo, mockDeleteFile)
	h := pengumuman_delete.NewDeletePengumumanHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengumuman/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengumuman", Value: "1"}}
	h.DeletePengumuman(w, req, ps)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
