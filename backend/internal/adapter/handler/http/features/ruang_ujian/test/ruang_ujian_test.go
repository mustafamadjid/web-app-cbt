package ruangujian_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	ruangujian_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/create"
	ruangujian_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/get"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ruangujian_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/create"
	ruangujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/get"
	ruangujianquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRuangUjianRepo struct {
	mock.Mock
}

func (m *MockRuangUjianRepo) GetRuangUjian(ctx context.Context, filter ruangujianquery.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]ruangujian.RuangUjian), args.Error(1)
}

func (m *MockRuangUjianRepo) GetRuangUjianById(ctx context.Context, idRuangan int) (ruangujian.RuangUjian, error) {
	args := m.Called(ctx, idRuangan)
	return args.Get(0).(ruangujian.RuangUjian), args.Error(1)
}

func (m *MockRuangUjianRepo) GetRuangUjianByKode(ctx context.Context, kodeRuang string) (ruangujian.RuangUjian, error) {
	args := m.Called(ctx, kodeRuang)
	return args.Get(0).(ruangujian.RuangUjian), args.Error(1)
}

func (m *MockRuangUjianRepo) ExistByKodeRuang(ctx context.Context, kodeRuang string) (bool, error) {
	args := m.Called(ctx, kodeRuang)
	return args.Bool(0), args.Error(1)
}

func (m *MockRuangUjianRepo) CreateRuangUjian(ctx context.Context, r ruangujian.RuangUjian) error {
	return m.Called(ctx, r).Error(0)
}

func (m *MockRuangUjianRepo) UpdateRuangUjian(ctx context.Context, idRuangan int, r updatepatch.UpdateRuangUjianPatch) error {
	return m.Called(ctx, idRuangan, r).Error(0)
}

func (m *MockRuangUjianRepo) DeleteRuangUjian(ctx context.Context, idRuangan int) error {
	return m.Called(ctx, idRuangan).Error(0)
}

func TestRuangUjianHandlers(t *testing.T) {
	mockRepo := new(MockRuangUjianRepo)
	createSvc := ruangujian_create_service.NewRuangUjianService(mockRepo)
	getSvc := ruangujian_get_service.NewGetRuangUjianService(mockRepo)

	createHandler := ruangujian_create.NewCreateRuangUjianHandler(createSvc)
	getHandler := ruangujian_get.NewGetRuangUjianHandler(getSvc)

	t.Run("Create Ruang Ujian Success", func(t *testing.T) {
		reqBody := `{"nama_ruangan":"Ruangan 1", "kode_ruang":"R1"}`
		req := httptest.NewRequest(http.MethodPost, "/ruang-ujian", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockRepo.On("ExistByKodeRuang", mock.Anything, "R1").Return(false, nil).Once()
		mockRepo.On("CreateRuangUjian", mock.Anything, mock.Anything).Return(nil).Once()

		createHandler.CreateRuangUjian(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Ruang Ujian By ID Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ruang-ujian/1", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "IdRuangan", Value: "1"}}

		mockRepo.On("GetRuangUjianById", mock.Anything, 1).Return(ruangujian.RuangUjian{IdRuangan: 1, NamaRuangan: "Ruangan 1"}, nil).Once()

		getHandler.GetRuangUjianByID(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
