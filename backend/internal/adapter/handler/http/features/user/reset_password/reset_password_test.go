package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	reset_password "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/reset_password"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/reset_password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResetPasswordRepo struct{ mock.Mock }

func (m *MockResetPasswordRepo) ResetPassword(ctx context.Context, id user.ID, password string) error {
	return m.Called(ctx, id, password).Error(0)
}

type MockHasher struct{ mock.Mock }

func (m *MockHasher) GenerateHash(plain string) (string, error) {
	a := m.Called(plain); return a.String(0), a.Error(1)
}
func (m *MockHasher) ComparePaswordAndHashed(hash, plain string) bool {
	return m.Called(hash, plain).Bool(0)
}

func TestResetPassword_Success(t *testing.T) {
	mockRepo := new(MockResetPasswordRepo)
	mockHasher := new(MockHasher)
	svc := user_service.NewResetPasswordService(mockRepo, mockHasher)
	h := reset_password.NewResetPasswordHandler(svc)
	req := httptest.NewRequest(http.MethodPut, "/pengguna/1/reset-password", strings.NewReader(`{"password":"NewPass123!"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengguna", Value: "1"}}
	mockHasher.On("GenerateHash", "NewPass123!").Return("hashedpw", nil).Once()
	mockRepo.On("ResetPassword", mock.Anything, user.ID(1), "hashedpw").Return(nil).Once()
	h.ResetPasswordHandler(w, req, ps)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResetPassword_BadRequest_InvalidParam(t *testing.T) {
	mockRepo := new(MockResetPasswordRepo)
	mockHasher := new(MockHasher)
	svc := user_service.NewResetPasswordService(mockRepo, mockHasher)
	h := reset_password.NewResetPasswordHandler(svc)
	req := httptest.NewRequest(http.MethodPut, "/pengguna/abc/reset-password", strings.NewReader(`{"password":"NewPass123!"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengguna", Value: "abc"}}
	h.ResetPasswordHandler(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id pengguna")
}

func TestResetPassword_BadRequest_MissingContentType(t *testing.T) {
	mockRepo := new(MockResetPasswordRepo)
	mockHasher := new(MockHasher)
	svc := user_service.NewResetPasswordService(mockRepo, mockHasher)
	h := reset_password.NewResetPasswordHandler(svc)
	req := httptest.NewRequest(http.MethodPut, "/pengguna/1/reset-password", strings.NewReader(`{"password":"NewPass123!"}`))
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengguna", Value: "1"}}
	h.ResetPasswordHandler(w, req, ps)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "content type must be application/json")
}

func TestResetPassword_MethodNotAllowed(t *testing.T) {
	mockRepo := new(MockResetPasswordRepo)
	mockHasher := new(MockHasher)
	svc := user_service.NewResetPasswordService(mockRepo, mockHasher)
	h := reset_password.NewResetPasswordHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/pengguna/1/reset-password", nil)
	w := httptest.NewRecorder()
	ps := httprouter.Params{{Key: "idPengguna", Value: "1"}}
	h.ResetPasswordHandler(w, req, ps)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
