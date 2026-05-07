package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	sesi_delete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/delete"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesiquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	sesi_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/delete"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSesiRepo struct{ mock.Mock }

func (m *MockSesiRepo) GetSesi(ctx context.Context, f sesiquery.ListSesiFilter) ([]sesi.Sesi, error) {
	a := m.Called(ctx, f); return a.Get(0).([]sesi.Sesi), a.Error(1)
}
func (m *MockSesiRepo) GetSesiById(ctx context.Context, id int) (sesi.Sesi, error) {
	a := m.Called(ctx, id); return a.Get(0).(sesi.Sesi), a.Error(1)
}
func (m *MockSesiRepo) GetSesiByKode(ctx context.Context, k string) (sesi.Sesi, error) {
	a := m.Called(ctx, k); return a.Get(0).(sesi.Sesi), a.Error(1)
}
func (m *MockSesiRepo) ExistByKodeSesi(ctx context.Context, k string) (bool, error) {
	a := m.Called(ctx, k); return a.Bool(0), a.Error(1)
}
func (m *MockSesiRepo) CreateSesi(ctx context.Context, s sesi.Sesi) error { return m.Called(ctx, s).Error(0) }
func (m *MockSesiRepo) UpdateSesi(ctx context.Context, id int, s updatepatch.UpdateSesiPatch) error {
	return m.Called(ctx, id, s).Error(0)
}
func (m *MockSesiRepo) DeleteSesi(ctx context.Context, id int) error { return m.Called(ctx, id).Error(0) }

func TestDeleteSesi_Success(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_delete_service.NewDeleteSesiService(mockRepo)
	h := sesi_delete.NewDeleteSesiHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/sesi/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "1"}}
	mockRepo.On("DeleteSesi", mock.Anything, 1).Return(nil).Once()
	h.DeleteSesi(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteSesi_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_delete_service.NewDeleteSesiService(mockRepo)
	h := sesi_delete.NewDeleteSesiHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/sesi/0", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "0"}}
	h.DeleteSesi(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id sesi")
}

func TestDeleteSesi_NotFound(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_delete_service.NewDeleteSesiService(mockRepo)
	h := sesi_delete.NewDeleteSesiHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/sesi/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "99"}}
	mockRepo.On("DeleteSesi", mock.Anything, 99).Return(coreerror.ErrNotFound).Once()
	h.DeleteSesi(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSesi_Restricted(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_delete_service.NewDeleteSesiService(mockRepo)
	h := sesi_delete.NewDeleteSesiHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/sesi/2", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idSesi", Value: "2"}}
	mockRepo.On("DeleteSesi", mock.Anything, 2).Return(coreerror.ErrDeleteRestricted).Once()
	h.DeleteSesi(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "DELETE_RESTRICTED")
}
