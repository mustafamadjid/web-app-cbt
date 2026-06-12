package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	kelas_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/get"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	kelasquery "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	kelas_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
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

func TestListKelas_Success(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)
	handler := kelas_get.NewGetKelasHandler(getSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas?q=A&tingkat_kelas=10&limit=20&offset=0", nil)
	w := httptest.NewRecorder()

	fullData := []kelas.FullKelasData{
		{
			ItemsTingkatKelas: []kelas.TingkatKelas{{IdTingkatKelas: 1, TingkatKelas: 10}},
			ItemsNamaKelas:    []kelas.NamaKelas{{IdNamaKelas: 1, IdTingkatKelas: 1, NamaKelas: "A"}},
		},
	}
	mockRepo.On("GetKelas", mock.Anything, mock.MatchedBy(func(filter kelasquery.ListKelasFilter) bool {
		return filter.Search == "A" && filter.TingkatKelas != nil && *filter.TingkatKelas == 10 && filter.Limit == 20 && filter.Offset == 0
	})).Return(fullData, nil).Once()

	handler.ListKelas(w, req, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "item_tingkat_kelas")
	mockRepo.AssertExpectations(t)
}

func TestListKelas_Success_DefaultParams(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)
	handler := kelas_get.NewGetKelasHandler(getSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas", nil)
	w := httptest.NewRecorder()

	mockRepo.On("GetKelas", mock.Anything, mock.Anything).Return([]kelas.FullKelasData{}, nil).Once()

	handler.ListKelas(w, req, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestListKelas_BadRequest_InvalidFilter(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)
	handler := kelas_get.NewGetKelasHandler(getSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas?tingkat_kelas=abc", nil)
	w := httptest.NewRecorder()

	handler.ListKelas(w, req, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "tingkat_kelas must be a number")
}

func TestGetKelasByID_Success(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)
	handler := kelas_get.NewGetKelasHandler(getSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas/1/2", nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "idTingkatKelas", Value: "1"}, {Key: "idNamaKelas", Value: "2"}}

	mockRepo.On("GetKelasById", mock.Anything, 1, 2).Return(kelas.KelasData{
		ItemsTingkatKelas: kelas.TingkatKelas{IdTingkatKelas: 1, TingkatKelas: 10},
		ItemsNamaKelas:    kelas.NamaKelas{IdNamaKelas: 2, IdTingkatKelas: 1, NamaKelas: "A"},
	}, nil).Once()

	handler.GetKelasByID(w, req, params)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "A")
	mockRepo.AssertExpectations(t)
}

func TestGetKelasByID_BadRequest_InvalidTingkatKelasParam(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)
	handler := kelas_get.NewGetKelasHandler(getSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas/x/2", nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "idTingkatKelas", Value: "x"}, {Key: "idNamaKelas", Value: "2"}}

	handler.GetKelasByID(w, req, params)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id_tingkat_kelas")
}

func TestGetKelasByID_BadRequest_InvalidNamaKelasParam(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)
	handler := kelas_get.NewGetKelasHandler(getSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas/1/abc", nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "idTingkatKelas", Value: "1"}, {Key: "idNamaKelas", Value: "abc"}}

	handler.GetKelasByID(w, req, params)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id_nama_kelas")
}

func TestGetKelasByID_NotFound(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)
	handler := kelas_get.NewGetKelasHandler(getSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas/1/99", nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "idTingkatKelas", Value: "1"}, {Key: "idNamaKelas", Value: "99"}}

	mockRepo.On("GetKelasById", mock.Anything, 1, 99).Return(kelas.KelasData{}, coreerror.ErrNotFound).Once()

	handler.GetKelasByID(w, req, params)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestListKelas_InternalServerError(t *testing.T) {
	mockRepo := new(MockKelasRepository)
	getSvc := kelas_get_service.NewGetKelasService(mockRepo)
	handler := kelas_get.NewGetKelasHandler(getSvc)

	req := httptest.NewRequest(http.MethodGet, "/kelas?limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	mockRepo.On("GetKelas", mock.Anything, mock.Anything).Return([]kelas.FullKelasData{}, coreerror.ErrDbError).Once()

	handler.ListKelas(w, req, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}
