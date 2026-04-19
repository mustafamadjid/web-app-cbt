package ujian_siswa_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	ujian_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/create"
	ujian_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/get"
	ujian_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/update"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	ujian_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/create"
	ujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/get"
	ujian_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUjianRepository struct {
	mock.Mock
}

func (m *MockUjianRepository) GetIdUjianByAttempt(ctx context.Context, idAttempt ujian.ID) (ujian.ID, error) {
	args := m.Called(ctx, idAttempt)
	return args.Get(0).(ujian.ID), args.Error(1)
}

func (m *MockUjianRepository) CreateUjian(ctx context.Context, u ujian.PenjadwalanUjian) error {
	return m.Called(ctx, u).Error(0)
}

func (m *MockUjianRepository) UpdateUjian(ctx context.Context, id ujian.ID, p updatepatch.UpdatePenjadwalanUjian) error {
	return m.Called(ctx, id, p).Error(0)
}

func (m *MockUjianRepository) DeleteUjian(ctx context.Context, id ujian.ID) error {
	return m.Called(ctx, id).Error(0)
}

type MockListUjianRepository struct {
	mock.Mock
}

func (m *MockListUjianRepository) GetAllUjian(ctx context.Context, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]ujian.ListUjian), args.Error(1)
}

func (m *MockListUjianRepository) GetUjianById(ctx context.Context, id ujian.ID) (ujian.ListUjian, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ujian.ListUjian), args.Error(1)
}

func (m *MockListUjianRepository) GetAllUjianSubmittedByIdSiswa(ctx context.Context, idSiswa int) ([]ujian.ListUjian, error) {
	args := m.Called(ctx, idSiswa)
	return args.Get(0).([]ujian.ListUjian), args.Error(1)
}

func TestAdminUjianHandlers(t *testing.T) {
	mockRepo := new(MockUjianRepository)
	mockListRepo := new(MockListUjianRepository)

	createSvc := ujian_create_service.NewCreateUjianService(mockRepo)
	getSvc := ujian_get_service.NewGetujianService(mockListRepo)
	updateSvc := ujian_update_service.NewUpdateUjianService(mockRepo)

	createHandler := ujian_create.NewCreateUjianHandler(createSvc)
	getHandler := ujian_get.NewGetUjianHandler(getSvc)
	updateHandler := ujian_update.NewUpdateUjianHandler(updateSvc)

	t.Run("Create Ujian Success", func(t *testing.T) {
		reqBody := `{
			"nama_ujian":"Ujian Akhir Semester",
			"id_bank_soal":1,
			"id_kelas":1,
			"id_nama_kelas":1,
			"id_guru":1,
			"id_sesi":1,
			"id_ruangan":1,
			"id_pengawas":1,
			"tanggal_ujian":"2023-10-10T00:00:00Z",
			"waktu_mulai":"2023-10-10T08:00:00Z",
			"waktu_selesai":"2023-10-10T10:00:00Z",
			"status_ujian":"BELUM_MULAI",
			"token":"ABCDEF"
		}`
		req := httptest.NewRequest(http.MethodPost, "/ujian", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockRepo.On("CreateUjian", mock.Anything, mock.Anything).Return(nil).Once()

		createHandler.CreateUjian(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create Ujian Conflict", func(t *testing.T) {
		reqBody := `{
			"nama_ujian":"Ujian Akhir Semester",
			"id_bank_soal":1,
			"id_kelas":1,
			"id_nama_kelas":1,
			"id_guru":1,
			"id_sesi":1,
			"id_ruangan":1,
			"id_pengawas":1,
			"tanggal_ujian":"2023-10-10T00:00:00Z",
			"waktu_mulai":"2023-10-10T08:00:00Z",
			"waktu_selesai":"2023-10-10T10:00:00Z",
			"status_ujian":"BELUM_MULAI",
			"token":"ABCDEF"
		}`
		req := httptest.NewRequest(http.MethodPost, "/ujian", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockRepo.On("CreateUjian", mock.Anything, mock.Anything).Return(coreerror.ErrConflict).Once()

		createHandler.CreateUjian(w, req, nil)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "jadwal ujian bentrok")
	})

	t.Run("Get Ujian By ID Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ujian/1", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idUjian", Value: "1"}}

		mockListRepo.On("GetUjianById", mock.Anything, ujian.ID(1)).Return(ujian.ListUjian{IdJadwalUjian: 1, NamaUjian: "Ujian 1"}, nil).Once()

		getHandler.GetUjianById(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update Ujian Conflict", func(t *testing.T) {
		reqBody := `{"id_sesi":1}`
		req := httptest.NewRequest(http.MethodPatch, "/ujian/detail/1", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idUjian", Value: "1"}}

		mockRepo.On("UpdateUjian", mock.Anything, ujian.ID(1), mock.Anything).Return(coreerror.ErrConflict).Once()

		updateHandler.UpdateUjian(w, req, params)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "jadwal ujian bentrok")
	})
}
