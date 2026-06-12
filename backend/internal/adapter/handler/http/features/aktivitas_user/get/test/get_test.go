package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	aktivitas_user_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/aktivitas_user/get"
	aktivitas_user "github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAktivitasUserRepo struct{ mock.Mock }

func (m *MockAktivitasUserRepo) CreateAktivitasUser(ctx context.Context, a aktivitas_user.AktivitasUser) error {
	return m.Called(ctx, a).Error(0)
}
func (m *MockAktivitasUserRepo) GetAktivitasUser(ctx context.Context) ([]aktivitas_user.AktivitasUser, error) {
	a := m.Called(ctx); return a.Get(0).([]aktivitas_user.AktivitasUser), a.Error(1)
}

func TestGetAktivitasUser_Success(t *testing.T) {
	mockRepo := new(MockAktivitasUserRepo)
	svc := aktivitas_user_service.NewAktivitasUserService(mockRepo)
	h := aktivitas_user_get.NewAktivitasUserHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/aktivitas-user", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetAktivitasUser", mock.Anything).Return([]aktivitas_user.AktivitasUser{
		{IdAktivitas: "act-1", IdPengguna: 1, Action: aktivitas_user.CREATE, Description: "Test"},
	}, nil).Once()
	h.GetAktivitasUser(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test")
}

func TestGetAktivitasUser_EmptyResult(t *testing.T) {
	mockRepo := new(MockAktivitasUserRepo)
	svc := aktivitas_user_service.NewAktivitasUserService(mockRepo)
	h := aktivitas_user_get.NewAktivitasUserHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/aktivitas-user", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetAktivitasUser", mock.Anything).Return([]aktivitas_user.AktivitasUser{}, nil).Once()
	h.GetAktivitasUser(w, req, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAktivitasUser_InternalServerError(t *testing.T) {
	mockRepo := new(MockAktivitasUserRepo)
	svc := aktivitas_user_service.NewAktivitasUserService(mockRepo)
	h := aktivitas_user_get.NewAktivitasUserHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/aktivitas-user", nil)
	w := httptest.NewRecorder()
	mockRepo.On("GetAktivitasUser", mock.Anything).Return([]aktivitas_user.AktivitasUser{}, coreerror.ErrDbError).Once()
	h.GetAktivitasUser(w, req, nil)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAktivitasUser_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockAktivitasUserRepo)
	svc := aktivitas_user_service.NewAktivitasUserService(mockRepo)
	h := aktivitas_user_get.NewAktivitasUserHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/aktivitas-user", nil)
	w := httptest.NewRecorder()
	h.GetAktivitasUser(w, req, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
