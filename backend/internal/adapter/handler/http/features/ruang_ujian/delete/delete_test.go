package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	ruangujian_delete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/delete"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ruangujianquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangujian_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/delete"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRuangUjianRepo struct{ mock.Mock }

func (m *MockRuangUjianRepo) GetRuangUjian(ctx context.Context, f ruangujianquery.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
	a := m.Called(ctx, f); return a.Get(0).([]ruangujian.RuangUjian), a.Error(1)
}
func (m *MockRuangUjianRepo) GetRuangUjianById(ctx context.Context, id int) (ruangujian.RuangUjian, error) {
	a := m.Called(ctx, id); return a.Get(0).(ruangujian.RuangUjian), a.Error(1)
}
func (m *MockRuangUjianRepo) GetRuangUjianByKode(ctx context.Context, k string) (ruangujian.RuangUjian, error) {
	a := m.Called(ctx, k); return a.Get(0).(ruangujian.RuangUjian), a.Error(1)
}
func (m *MockRuangUjianRepo) ExistByKodeRuang(ctx context.Context, k string) (bool, error) {
	a := m.Called(ctx, k); return a.Bool(0), a.Error(1)
}
func (m *MockRuangUjianRepo) CreateRuangUjian(ctx context.Context, r ruangujian.RuangUjian) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockRuangUjianRepo) UpdateRuangUjian(ctx context.Context, id int, r updatepatch.UpdateRuangUjianPatch) error {
	return m.Called(ctx, id, r).Error(0)
}
func (m *MockRuangUjianRepo) DeleteRuangUjian(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestDeleteRuangUjian_Success(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_delete_service.NewDeleteRuangUjianService(mockRepo)
	h := ruangujian_delete.NewDeleteRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/ruang-ujian/1", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idRuangan", Value: "1"}}
	mockRepo.On("DeleteRuangUjian", mock.Anything, 1).Return(nil).Once()
	h.DeleteRuangUjian(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteRuangUjian_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_delete_service.NewDeleteRuangUjianService(mockRepo)
	h := ruangujian_delete.NewDeleteRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/ruang-ujian/abc", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idRuangan", Value: "abc"}}
	h.DeleteRuangUjian(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id ruangan")
}

func TestDeleteRuangUjian_NotFound(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_delete_service.NewDeleteRuangUjianService(mockRepo)
	h := ruangujian_delete.NewDeleteRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/ruang-ujian/99", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idRuangan", Value: "99"}}
	mockRepo.On("DeleteRuangUjian", mock.Anything, 99).Return(coreerror.ErrNotFound).Once()
	h.DeleteRuangUjian(w, req, ps)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteRuangUjian_Restricted(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_delete_service.NewDeleteRuangUjianService(mockRepo)
	h := ruangujian_delete.NewDeleteRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/ruang-ujian/2", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idRuangan", Value: "2"}}
	mockRepo.On("DeleteRuangUjian", mock.Anything, 2).Return(coreerror.ErrDeleteRestricted).Once()
	h.DeleteRuangUjian(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "DELETE_RESTRICTED")
}
