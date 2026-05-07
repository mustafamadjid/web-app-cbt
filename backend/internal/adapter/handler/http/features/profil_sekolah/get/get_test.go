package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	profil_sekolah_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/get"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	profil_sekolah "github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	profil_sekolah_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProfilSekolahRepo struct{ mock.Mock }

func (m *MockProfilSekolahRepo) UpdateProfilSekolah(ctx context.Context, id profil_sekolah.IDProfil, p profil_sekolah.ProfilSekolah) error {
	return m.Called(ctx, id, p).Error(0)
}
func (m *MockProfilSekolahRepo) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	a := m.Called(ctx); return a.Get(0).(profil_sekolah.ProfilSekolah), a.Error(1)
}

func TestGetProfilSekolah_Success(t *testing.T) {
	mockRepo := new(MockProfilSekolahRepo)
	svc := profil_sekolah_get_service.NewGetProfilSekolahService(mockRepo)
	h := profil_sekolah_get.NewGetProfilSekolahHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/profil-sekolah", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetProfilSekolah", mock.Anything).Return(profil_sekolah.ProfilSekolah{NamaSekolah: "SMA Test"}, nil).Once()
	h.GetProfilSekolah(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SMA Test")
}

func TestGetProfilSekolah_NotFound(t *testing.T) {
	mockRepo := new(MockProfilSekolahRepo)
	svc := profil_sekolah_get_service.NewGetProfilSekolahService(mockRepo)
	h := profil_sekolah_get.NewGetProfilSekolahHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/profil-sekolah", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetProfilSekolah", mock.Anything).Return(profil_sekolah.ProfilSekolah{}, coreerror.ErrNotFound).Once()
	h.GetProfilSekolah(w, req, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetProfilSekolah_InternalServerError(t *testing.T) {
	mockRepo := new(MockProfilSekolahRepo)
	svc := profil_sekolah_get_service.NewGetProfilSekolahService(mockRepo)
	h := profil_sekolah_get.NewGetProfilSekolahHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/profil-sekolah", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetProfilSekolah", mock.Anything).Return(profil_sekolah.ProfilSekolah{}, coreerror.ErrDbError).Once()
	h.GetProfilSekolah(w, req, nil)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetProfilSekolah_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockProfilSekolahRepo)
	svc := profil_sekolah_get_service.NewGetProfilSekolahService(mockRepo)
	h := profil_sekolah_get.NewGetProfilSekolahHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/profil-sekolah", nil)
	w := httptest.NewRecorder()
	h.GetProfilSekolah(w, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
