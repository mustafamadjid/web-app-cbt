package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	kelas_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/create"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	kelasquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	kelas_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKelasRepository struct {
	mock.Mock
}

func (m *MockKelasRepository) GetKelas(ctx context.Context, filter kelasquery.ListKelasFilter) ([]kelas.FullKelasData, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]kelas.FullKelasData), args.Error(1)
}

func (m *MockKelasRepository) GetKelasById(ctx context.Context, idTingkatKelas int, idNamaKelas int) (kelas.KelasData, error) {
	args := m.Called(ctx, idTingkatKelas, idNamaKelas)
	return args.Get(0).(kelas.KelasData), args.Error(1)
}

func (m *MockKelasRepository) CreateTingkatKelas(ctx context.Context, tingkatKelas int) error {
	return m.Called(ctx, tingkatKelas).Error(0)
}

func (m *MockKelasRepository) CreateNamaKelas(ctx context.Context, namaKelas kelas.NamaKelas) error {
	return m.Called(ctx, namaKelas).Error(0)
}

func (m *MockKelasRepository) UpdateNamaKelas(ctx context.Context, idNamaKelas int, dataUpdate updatepatch.NamaKelasPatch) error {
	return m.Called(ctx, idNamaKelas, dataUpdate).Error(0)
}

func (m *MockKelasRepository) DeleteNamaKelas(ctx context.Context, idNamaKelas int) error {
	return m.Called(ctx, idNamaKelas).Error(0)
}

func (m *MockKelasRepository) ExistTingkatKelas(ctx context.Context, tingkatKelas int) (bool, error) {
	args := m.Called(ctx, tingkatKelas)
	return args.Bool(0), args.Error(1)
}

func (m *MockKelasRepository) ExistNamaKelas(ctx context.Context, namaKelas string) (bool, error) {
	args := m.Called(ctx, namaKelas)
	return args.Bool(0), args.Error(1)
}

func withAdminActor(ctx context.Context) context.Context {
	return middleware.WithActor(ctx, user.Actor{IdPengguna: 1, Role: user.ADMIN})
}

func TestCreateNamaKelas_Success(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	reqBody := `{"id_tingkat_kelas":1, "nama_kelas":"Kelas A"}`
	req := httptest.NewRequest(http.MethodPost, "/kelas", strings.NewReader(reqBody)).WithContext(withAdminActor(context.Background()))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("ExistNamaKelas", mock.Anything, "Kelas A").Return(false, nil).Once()
	mockRepo.On("CreateNamaKelas", mock.Anything, mock.Anything).Return(nil).Once()

	handler.CreateNamaKelas(w, req, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateTingkatKelas_Success(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodPost, "/kelas/tingkat", strings.NewReader(`{"tingkat_kelas":10}`)).WithContext(withAdminActor(context.Background()))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("ExistTingkatKelas", mock.Anything, 10).Return(false, nil).Once()
	mockRepo.On("CreateTingkatKelas", mock.Anything, 10).Return(nil).Once()

	handler.CreateTingkatKelas(w, req, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateNamaKelas_BadRequest_MissingContentType(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodPost, "/kelas", strings.NewReader(`{"id_tingkat_kelas":1, "nama_kelas":"A"}`))
	w := httptest.NewRecorder()

	handler.CreateNamaKelas(w, req, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "content type must be application/json")
}

func TestCreateNamaKelas_BadRequest_InvalidJSON(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodPost, "/kelas", strings.NewReader(`{"id_tingkat_kelas":1`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateNamaKelas(w, req, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestCreateNamaKelas_BadRequest_RequiredFieldMissing(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodPost, "/kelas", strings.NewReader(`{"id_tingkat_kelas":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateNamaKelas(w, req, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "nama kelas is required")
}

func TestCreateNamaKelas_BadRequest_AlreadyExists(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodPost, "/kelas", strings.NewReader(`{"id_tingkat_kelas":1, "nama_kelas":"Kelas A"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("ExistNamaKelas", mock.Anything, "Kelas A").Return(true, nil).Once()

	handler.CreateNamaKelas(w, req, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exist")
}

func TestCreateTingkatKelas_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas/tingkat", nil)
	w := httptest.NewRecorder()

	handler.CreateTingkatKelas(w, req, nil)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestCreateTingkatKelas_BadRequest_Empty(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodPost, "/kelas/tingkat", strings.NewReader(`{"tingkat_kelas":0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTingkatKelas(w, req, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "tingkat kelas is required")
}

func TestCreateTingkatKelas_BadRequest_AlreadyExists(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodPost, "/kelas/tingkat", strings.NewReader(`{"tingkat_kelas":10}`)).WithContext(withAdminActor(context.Background()))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("ExistTingkatKelas", mock.Anything, 10).Return(true, nil).Once()

	handler.CreateTingkatKelas(w, req, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exist")
}

func TestCreateNamaKelas_BadRequest_IdTingkatKelasInvalid(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	req := httptest.NewRequest(http.MethodPost, "/kelas", strings.NewReader(`{"id_tingkat_kelas":0, "nama_kelas":"A"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateNamaKelas(w, req, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "id tingkat kelas is required")
}

func TestCreateNamaKelas_InternalServerError_RepoFails(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	handler := kelas_create.NewCreateKelasHandler(createSvc)

	reqBody := `{"id_tingkat_kelas":1, "nama_kelas":"Kelas X"}`
	req := httptest.NewRequest(http.MethodPost, "/kelas", strings.NewReader(reqBody)).WithContext(withAdminActor(context.Background()))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("ExistNamaKelas", mock.Anything, "Kelas X").Return(false, nil).Once()
	mockRepo.On("CreateNamaKelas", mock.Anything, mock.Anything).Return(coreerror.ErrDbError).Once()

	handler.CreateNamaKelas(w, req, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
