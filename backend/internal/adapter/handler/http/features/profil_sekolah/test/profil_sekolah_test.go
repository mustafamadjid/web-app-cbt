package profil_sekolah_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/get"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	profil_sekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProfilSekolahRepository struct {
	mock.Mock
}

func (m *MockProfilSekolahRepository) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	args := m.Called(ctx)
	return args.Get(0).(profil_sekolah.ProfilSekolah), args.Error(1)
}

func (m *MockProfilSekolahRepository) UpdateProfilSekolah(ctx context.Context, id profil_sekolah.IDProfil, p profil_sekolah.ProfilSekolah) error {
	return m.Called(ctx, id, p).Error(0)
}

func TestProfilSekolahGetHandler(t *testing.T) {
	mockRepo := new(MockProfilSekolahRepository)
	svc := profil_sekolah_service.NewGetProfilSekolahService(mockRepo)
	handler := httpx.NewGetProfilSekolahHandler(svc)

	t.Run("Get Profil Sekolah Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/profil-sekolah", nil)
		w := httptest.NewRecorder()

		mockRepo.On("GetProfilSekolah", mock.Anything).Return(profil_sekolah.ProfilSekolah{NamaSekolah: "Test School"}, nil).Once()

		handler.GetProfilSekolah(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test School")
	})
}
