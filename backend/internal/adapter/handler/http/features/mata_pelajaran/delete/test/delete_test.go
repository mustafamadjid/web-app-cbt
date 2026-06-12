package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	matapelajaran_delete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/delete"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	matapelajaranquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	matapelajaran_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/delete"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMapelRepo struct{ mock.Mock }

func (m *MockMapelRepo) GetMapel(ctx context.Context, f matapelajaranquery.ListMapelFilter) ([]matapelajaran.MataPelajaran, error) {
	a := m.Called(ctx, f); return a.Get(0).([]matapelajaran.MataPelajaran), a.Error(1)
}
func (m *MockMapelRepo) GetMapelById(ctx context.Context, id int) (matapelajaran.MataPelajaran, error) {
	a := m.Called(ctx, id); return a.Get(0).(matapelajaran.MataPelajaran), a.Error(1)
}
func (m *MockMapelRepo) CreateMapel(ctx context.Context, s matapelajaran.MataPelajaran) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockMapelRepo) UpdateMapel(ctx context.Context, id int, s updatepatch.UpdateMapelPatch) error {
	return m.Called(ctx, id, s).Error(0)
}
func (m *MockMapelRepo) DeleteMapel(ctx context.Context, id int) error { return m.Called(ctx, id).Error(0) }
func (m *MockMapelRepo) ExistKodeMapel(ctx context.Context, k string) (bool, error) {
	a := m.Called(ctx, k); return a.Bool(0), a.Error(1)
}

func TestDeleteMapel_Success(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_delete_service.NewDeleteMapelService(mockRepo)
	h := matapelajaran_delete.NewDeleteMapelHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/mata-pelajaran/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "1"}}
	mockRepo.On("DeleteMapel", mock.Anything, 1).Return(nil).Once()
	h.DeleteMapel(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteMapel_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_delete_service.NewDeleteMapelService(mockRepo)
	h := matapelajaran_delete.NewDeleteMapelHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/mata-pelajaran/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "abc"}}
	h.DeleteMapel(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id mapel")
}

func TestDeleteMapel_NotFound(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_delete_service.NewDeleteMapelService(mockRepo)
	h := matapelajaran_delete.NewDeleteMapelHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/mata-pelajaran/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "99"}}
	mockRepo.On("DeleteMapel", mock.Anything, 99).Return(coreerror.ErrNotFound).Once()
	h.DeleteMapel(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteMapel_Restricted(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_delete_service.NewDeleteMapelService(mockRepo)
	h := matapelajaran_delete.NewDeleteMapelHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/mata-pelajaran/2", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idMapel", Value: "2"}}
	mockRepo.On("DeleteMapel", mock.Anything, 2).Return(coreerror.ErrDeleteRestricted).Once()
	h.DeleteMapel(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "DELETE_RESTRICTED")
}
