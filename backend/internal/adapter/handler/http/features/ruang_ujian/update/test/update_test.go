package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	ruangujian_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/update"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ruangujianquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangujian_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/update"
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

func TestUpdateRuangUjian_Success(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_update_service.NewUpdateRuangUjianService(mockRepo)
	h := ruangujian_update.NewUpdateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/ruang-ujian/1", strings.NewReader(`{"nama_ruangan":"Lab 2"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idRuangan", Value: "1"}}
	mockRepo.On("UpdateRuangUjian", mock.Anything, 1, mock.MatchedBy(func(p updatepatch.UpdateRuangUjianPatch) bool {
		return p.NamaRuang != nil && *p.NamaRuang == "Lab 2"
	})).Return(nil).Once()
	h.UpdateRuangUjian(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateRuangUjian_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_update_service.NewUpdateRuangUjianService(mockRepo)
	h := ruangujian_update.NewUpdateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/ruang-ujian/abc", strings.NewReader(`{"nama_ruangan":"Lab 2"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idRuangan", Value: "abc"}}
	h.UpdateRuangUjian(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id ruangan")
}

func TestUpdateRuangUjian_BadRequest_KodeExists(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_update_service.NewUpdateRuangUjianService(mockRepo)
	h := ruangujian_update.NewUpdateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/ruang-ujian/1", strings.NewReader(`{"kode_ruang":"R2"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idRuangan", Value: "1"}}
	mockRepo.On("ExistByKodeRuang", mock.Anything, "R2").Return(true, nil).Once()
	h.UpdateRuangUjian(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exist")
}

func TestUpdateRuangUjian_InternalServerError(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_update_service.NewUpdateRuangUjianService(mockRepo)
	h := ruangujian_update.NewUpdateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPatch, "/ruang-ujian/99", strings.NewReader(`{"nama_ruangan":"Lab 2"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idRuangan", Value: "99"}}
	mockRepo.On("UpdateRuangUjian", mock.Anything, 99, mock.Anything).Return(coreerror.ErrNotFound).Once()
	h.UpdateRuangUjian(w, req, ps)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
