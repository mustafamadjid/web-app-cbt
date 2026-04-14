package aktivitas_user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/aktivitas_user/get"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAktivitasUserRepository struct {
	mock.Mock
}

func (m *MockAktivitasUserRepository) GetAktivitasUser(ctx context.Context) ([]aktivitas_user.AktivitasUser, error) {
	args := m.Called(ctx)
	return args.Get(0).([]aktivitas_user.AktivitasUser), args.Error(1)
}

func (m *MockAktivitasUserRepository) CreateAktivitasUser(ctx context.Context, a aktivitas_user.AktivitasUser) error {
	return m.Called(ctx, a).Error(0)
}

func TestAktivitasUserHandlers(t *testing.T) {
	mockRepo := new(MockAktivitasUserRepository)
	svc := aktivitas_user_service.NewAktivitasUserService(mockRepo)
	handler := httpx.NewAktivitasUserHandler(svc)

	t.Run("Get Aktivitas User Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/aktivitas-user", nil)
		w := httptest.NewRecorder()

		mockRepo.On("GetAktivitasUser", mock.Anything).Return([]aktivitas_user.AktivitasUser{{Description: "Test Description"}}, nil).Once()

		handler.GetAktivitasUser(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Description")
	})
}
