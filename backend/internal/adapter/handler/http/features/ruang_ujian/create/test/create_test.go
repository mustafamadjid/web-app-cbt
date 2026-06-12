package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ruangujian_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/create"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ruangujianquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangujian_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/create"
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

func TestCreateRuangUjian_Success(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_create_service.NewRuangUjianService(mockRepo)
	h := ruangujian_create.NewCreateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/ruang-ujian", strings.NewReader(`{"nama_ruangan":"Ruangan 1", "kode_ruang":"R1"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockRepo.On("ExistByKodeRuang", mock.Anything, "R1").Return(false, nil).Once()
	mockRepo.On("CreateRuangUjian", mock.Anything, mock.Anything).Return(nil).Once()
	h.CreateRuangUjian(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRuangUjian_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_create_service.NewRuangUjianService(mockRepo)
	h := ruangujian_create.NewCreateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/ruang-ujian", nil)
	w := httptest.NewRecorder()
	h.CreateRuangUjian(w, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestCreateRuangUjian_BadRequest_MissingContentType(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_create_service.NewRuangUjianService(mockRepo)
	h := ruangujian_create.NewCreateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/ruang-ujian", strings.NewReader(`{"nama_ruangan":"Ruangan 1", "kode_ruang":"R1"}`))
	w := httptest.NewRecorder()
	h.CreateRuangUjian(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "content type must be application/json")
}

func TestCreateRuangUjian_BadRequest_InvalidBody(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_create_service.NewRuangUjianService(mockRepo)
	h := ruangujian_create.NewCreateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/ruang-ujian", strings.NewReader(`{"nama_ruangan":"Ruangan 1"`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateRuangUjian(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestCreateRuangUjian_BadRequest_RequiredFieldMissing(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_create_service.NewRuangUjianService(mockRepo)
	h := ruangujian_create.NewCreateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/ruang-ujian", strings.NewReader(`{"nama_ruangan":"Ruangan 1"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateRuangUjian(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestCreateRuangUjian_BadRequest_KodeExists(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	svc := ruangujian_create_service.NewRuangUjianService(mockRepo)
	h := ruangujian_create.NewCreateRuangUjianHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/ruang-ujian", strings.NewReader(`{"nama_ruangan":"Ruangan 1", "kode_ruang":"R1"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockRepo.On("ExistByKodeRuang", mock.Anything, "R1").Return(true, nil).Once()
	h.CreateRuangUjian(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exist")
}
