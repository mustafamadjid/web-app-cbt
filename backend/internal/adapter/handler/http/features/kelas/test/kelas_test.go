package kelas_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	kelas_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/create"
	kelas_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/get"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	kelas_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	kelas_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
	kelasquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
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

func TestKelasHandlers(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	createSvc := kelas_create_service.NewCreateKelasService(mockRepo)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)

	createHandler := kelas_create.NewCreateKelasHandler(createSvc)
	getHandler := kelas_get.NewGetKelasHandler(getSvc)

	t.Run("Create Nama Kelas Success", func(t *testing.T) {
		reqBody := `{"id_tingkat_kelas":1, "nama_kelas":"Kelas A"}`
		actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
		ctx := middleware.WithActor(context.Background(), actor)

		req := httptest.NewRequest(http.MethodPost, "/kelas", strings.NewReader(reqBody)).WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockRepo.On("ExistTingkatKelas", mock.Anything, 1).Return(true, nil).Once()
		mockRepo.On("ExistNamaKelas", mock.Anything, "Kelas A").Return(false, nil).Once()
		mockRepo.On("CreateNamaKelas", mock.Anything, mock.Anything).Return(nil).Once()

		createHandler.CreateNamaKelas(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Kelas List Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/kelas", nil)
		w := httptest.NewRecorder()

		fullData := []kelas.FullKelasData{
			{
				ItemsTingkatKelas: []kelas.TingkatKelas{{IdTingkatKelas: 1, TingkatKelas: 10}},
				ItemsNamaKelas:    []kelas.NamaKelas{{IdNamaKelas: 1, IdTingkatKelas: 1, NamaKelas: "A"}},
			},
		}
		mockRepo.On("GetKelas", mock.Anything, mock.Anything).Return(fullData, nil).Once()

		getHandler.ListKelas(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
