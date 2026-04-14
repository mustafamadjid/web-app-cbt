package matapelajaran_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	matapelajaran_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/create"
	matapelajaran_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/get"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	matapelajaran_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/create"
	matapelajaran_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/get"
	matapelajaranquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMataPelajaranRepository struct {
	mock.Mock
}

func (m *MockMataPelajaranRepository) GetMapel(ctx context.Context, filter matapelajaranquery.ListMapelFilter) ([]matapelajaran.MataPelajaran, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]matapelajaran.MataPelajaran), args.Error(1)
}

func (m *MockMataPelajaranRepository) GetMapelById(ctx context.Context, idMapel int) (matapelajaran.MataPelajaran, error) {
	args := m.Called(ctx, idMapel)
	return args.Get(0).(matapelajaran.MataPelajaran), args.Error(1)
}

func (m *MockMataPelajaranRepository) CreateMapel(ctx context.Context, s matapelajaran.MataPelajaran) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockMataPelajaranRepository) UpdateMapel(ctx context.Context, idMapel int, s updatepatch.UpdateMapelPatch) error {
	return m.Called(ctx, idMapel, s).Error(0)
}

func (m *MockMataPelajaranRepository) DeleteMapel(ctx context.Context, idMapel int) error {
	return m.Called(ctx, idMapel).Error(0)
}

func (m *MockMataPelajaranRepository) ExistKodeMapel(ctx context.Context, kodeMapel string) (bool, error) {
	args := m.Called(ctx, kodeMapel)
	return args.Bool(0), args.Error(1)
}

func TestMataPelajaranHandlers(t *testing.T) {
	mockRepo := new(MockMataPelajaranRepository)
	createSvc := matapelajaran_create_service.NewMapelService(mockRepo)
	getSvc := matapelajaran_get_service.NewGetMapelService(mockRepo)

	createHandler := matapelajaran_create.NewCreateMapelHandler(createSvc)
	getHandler := matapelajaran_get.NewGetMapelHandler(getSvc)

	t.Run("Create Mata Pelajaran Success", func(t *testing.T) {
		reqBody := `{"id_kelas":1, "nama_mapel":"Matematika", "kode_mapel":"MTK", "deskripsi":"Matematika Dasar"}`
		req := httptest.NewRequest(http.MethodPost, "/mata-pelajaran", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockRepo.On("ExistKodeMapel", mock.Anything, "MTK").Return(false, nil).Once()
		mockRepo.On("CreateMapel", mock.Anything, mock.Anything).Return(nil).Once()

		createHandler.CreateMapel(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Mata Pelajaran By ID Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/mata-pelajaran/1", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idMapel", Value: "1"}}

		mockRepo.On("GetMapelById", mock.Anything, 1).Return(matapelajaran.MataPelajaran{IdMapel: 1, NamaMapel: "Matematika"}, nil).Once()

		getHandler.GetMapelByID(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
