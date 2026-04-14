package pengumuman_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	pengumuman_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/create"
	pengumuman_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/get"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/create"
	pengumuman_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPengumumanRepo struct {
	mock.Mock
}

func (m *MockPengumumanRepo) GetPengumumanActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	args := m.Called(ctx)
	return args.Get(0).([]pengumuman.Pengumuman), args.Error(1)
}

func (m *MockPengumumanRepo) GetPengumumanNonActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	args := m.Called(ctx)
	return args.Get(0).([]pengumuman.Pengumuman), args.Error(1)
}

func (m *MockPengumumanRepo) GetPengumumanIncoming(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	args := m.Called(ctx)
	return args.Get(0).([]pengumuman.Pengumuman), args.Error(1)
}

func (m *MockPengumumanRepo) GetPengumumanById(ctx context.Context, id pengumuman.ID) (pengumuman.Pengumuman, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(pengumuman.Pengumuman), args.Error(1)
}

func (m *MockPengumumanRepo) CreatePengumuman(ctx context.Context, p pengumuman.Pengumuman) error {
	return m.Called(ctx, p).Error(0)
}

func (m *MockPengumumanRepo) UpdatePengumuman(ctx context.Context, id pengumuman.ID, u updatepatch.PengumumanUpdatePatch) error {
	return m.Called(ctx, id, u).Error(0)
}

func (m *MockPengumumanRepo) DeletePengumuman(ctx context.Context, id pengumuman.ID) error {
	return m.Called(ctx, id).Error(0)
}

func TestPengumumanHandlers(t *testing.T) {
	mockRepo := new(MockPengumumanRepo)
	createSvc := pengumuman_create_service.NewCreatePengumumanRepo(mockRepo)
	getSvc := pengumuman_get_service.NewGetPengumumanService(mockRepo)

	// Mock DocumentStore
	docStore := &httphelper.DocumentStore{}

	createHandler := pengumuman_create.NewCreatePengumumanHandler(createSvc, *docStore)
	getHandler := pengumuman_get.NewGetPengumumanHandler(getSvc)

	t.Run("Create Pengumuman Success", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("judul_pengumuman", "Pengumuman 1")
		_ = writer.WriteField("isi_pengumuman", "Konten 1")
		_ = writer.WriteField("tanggal_rilis_pengumuman", "2023-10-10")
		_ = writer.WriteField("tanggal_selesai_pengumuman", "2023-10-20")
		_ = writer.Close()

		actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
		ctx := middleware.WithActor(context.Background(), actor)

		req := httptest.NewRequest(http.MethodPost, "/pengumuman", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		mockRepo.On("CreatePengumuman", mock.Anything, mock.Anything).Return(nil).Once()

		createHandler.CreatePengumuman(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Active Pengumuman Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pengumuman/active", nil)
		w := httptest.NewRecorder()

		mockRepo.On("GetPengumumanActive", mock.Anything).Return([]pengumuman.Pengumuman{{JudulPengumuman: "Active"}}, nil).Once()

		getHandler.GetPengumumanActive(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
