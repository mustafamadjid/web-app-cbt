package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sesi_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/create"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesiquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	sesi_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
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

func TestCreateSesi_Success(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_create_service.NewCreateSesiService(mockRepo)
	h := sesi_create.NewCreateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/sesi", strings.NewReader(`{"nama_sesi":"Sesi 1", "kode_sesi":"SESI1"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockRepo.On("ExistByKodeSesi", mock.Anything, "SESI1").Return(false, nil).Once()
	mockRepo.On("CreateSesi", mock.Anything, mock.Anything).Return(nil).Once()
	h.CreateSesiHandler(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateSesi_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_create_service.NewCreateSesiService(mockRepo)
	h := sesi_create.NewCreateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/sesi", nil)
	w := httptest.NewRecorder()
	h.CreateSesiHandler(w, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestCreateSesi_BadRequest_MissingContentType(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_create_service.NewCreateSesiService(mockRepo)
	h := sesi_create.NewCreateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/sesi", strings.NewReader(`{"nama_sesi":"Sesi 1"}`))
	w := httptest.NewRecorder()
	h.CreateSesiHandler(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "content type must be application/json")
}

func TestCreateSesi_BadRequest_InvalidBody(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_create_service.NewCreateSesiService(mockRepo)
	h := sesi_create.NewCreateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/sesi", strings.NewReader(`{"nama_sesi":"Sesi 1"`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateSesiHandler(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestCreateSesi_BadRequest_RequiredFieldMissing(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_create_service.NewCreateSesiService(mockRepo)
	h := sesi_create.NewCreateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/sesi", strings.NewReader(`{"nama_sesi":"Sesi 1"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateSesiHandler(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "kode sesi is required")
}

func TestCreateSesi_BadRequest_KodeExists(t *testing.T) {
	mockRepo := new(MockSesiRepo)
	svc := sesi_create_service.NewCreateSesiService(mockRepo)
	h := sesi_create.NewCreateSesiHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/sesi", strings.NewReader(`{"nama_sesi":"Sesi 1", "kode_sesi":"SESI1"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockRepo.On("ExistByKodeSesi", mock.Anything, "SESI1").Return(true, nil).Once()
	h.CreateSesiHandler(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exist")
}
