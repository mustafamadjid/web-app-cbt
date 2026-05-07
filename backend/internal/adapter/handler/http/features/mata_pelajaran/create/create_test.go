package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	matapelajaran_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/create"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	matapelajaranquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	matapelajaran_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/create"
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
func (m *MockMapelRepo) DeleteMapel(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockMapelRepo) ExistKodeMapel(ctx context.Context, kode string) (bool, error) {
	a := m.Called(ctx, kode); return a.Bool(0), a.Error(1)
}

func TestCreateMapel_Success(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_create_service.NewMapelService(mockRepo)
	h := matapelajaran_create.NewCreateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/mata-pelajaran", strings.NewReader(`{"id_kelas":1, "nama_mapel":"Matematika", "kode_mapel":"MTK", "deskripsi":"Matematika Dasar"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockRepo.On("ExistKodeMapel", mock.Anything, "MTK").Return(false, nil).Once()
	mockRepo.On("CreateMapel", mock.Anything, mock.Anything).Return(nil).Once()
	h.CreateMapel(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateMapel_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_create_service.NewMapelService(mockRepo)
	h := matapelajaran_create.NewCreateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/mata-pelajaran", nil)
	w := httptest.NewRecorder()
	h.CreateMapel(w, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestCreateMapel_BadRequest_MissingContentType(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_create_service.NewMapelService(mockRepo)
	h := matapelajaran_create.NewCreateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/mata-pelajaran", strings.NewReader(`{"id_kelas":1}`))
	w := httptest.NewRecorder()
	h.CreateMapel(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "content type must be application/json")
}

func TestCreateMapel_BadRequest_InvalidBody(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_create_service.NewMapelService(mockRepo)
	h := matapelajaran_create.NewCreateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/mata-pelajaran", strings.NewReader(`{"id_kelas":1`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateMapel(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestCreateMapel_BadRequest_RequiredFieldMissing(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_create_service.NewMapelService(mockRepo)
	h := matapelajaran_create.NewCreateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/mata-pelajaran", strings.NewReader(`{"id_kelas":1, "kode_mapel":"MTK", "deskripsi":"Matematika Dasar"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateMapel(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "nama mapel is required")
}

func TestCreateMapel_BadRequest_KodeExists(t *testing.T) {
	mockRepo := new(MockMapelRepo)
	svc := matapelajaran_create_service.NewMapelService(mockRepo)
	h := matapelajaran_create.NewCreateMapelHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/mata-pelajaran", strings.NewReader(`{"id_kelas":1, "nama_mapel":"Matematika", "kode_mapel":"MTK", "deskripsi":"Matematika Dasar"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockRepo.On("ExistKodeMapel", mock.Anything, "MTK").Return(true, nil).Once()
	h.CreateMapel(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exist")
}
