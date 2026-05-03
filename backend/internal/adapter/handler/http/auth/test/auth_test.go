package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/auth"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUsecase is a mock implementation of auth_port_in.AuthUsecase
type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Login(ctx context.Context, cmd auth_service.LoginCmd) (auth_service.LoginRes, error) {
	args := m.Called(ctx, cmd)
	return args.Get(0).(auth_service.LoginRes), args.Error(1)
}

func (m *MockAuthUsecase) Logout(ctx context.Context, refreshtoken string, now time.Time) error {
	args := m.Called(ctx, refreshtoken, now)
	return args.Error(0)
}

func (m *MockAuthUsecase) RefreshAccessToken(ctx context.Context, refreshToken string, accessTTL time.Duration) (string, error) {
	args := m.Called(ctx, refreshToken, accessTTL)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUsecase) AdminRevokingSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

// MockAccessTokenService is a mock implementation of out.AccessTokenService
type MockAccessTokenService struct {
	mock.Mock
}

func (m *MockAccessTokenService) GenerateAccessToken(idPengguna user.ID, role user.Role, username string, tokenDuration time.Duration) (string, error) {
	args := m.Called(idPengguna, role, username, tokenDuration)
	return args.String(0), args.Error(1)
}

func (m *MockAccessTokenService) VerifyAccessToken(token string, now time.Time) (user.ID, user.Role, string, error) {
	args := m.Called(token, now)
	return args.Get(0).(user.ID), args.Get(1).(user.Role), args.String(2), args.Error(3)
}

func TestLogin(t *testing.T) {
	mockSvc := new(MockAuthUsecase)
	cookieCfg := cookie.CookieConfig{
		AccessName:  "access_token",
		RefreshName: "refresh_token",
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
	}

	handler := httpx.NewAuthHandler(mockSvc, cookieCfg, time.Hour, time.Hour*24, nil, nil)

	t.Run("Success Login", func(t *testing.T) {
		reqBody := httpx.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSvc.On("Login", mock.Anything, auth_service.LoginCmd{
			Username: "testuser",
			Password: "password123",
		}).Return(auth_service.LoginRes{
			Username:     "testuser",
			AccessToken:  "access",
			RefreshToken: "refresh",
			IdPengguna:   1,
		}, nil).Once()

		handler.Login(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "success", resp["message"])

		// Check cookies
		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 2)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		reqBody := httpx.LoginRequest{
			Username: "wronguser",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSvc.On("Login", mock.Anything, auth_service.LoginCmd{
			Username: "wronguser",
			Password: "wrongpassword",
		}).Return(auth_service.LoginRes{}, coreerror.ErrInvalidCreds).Once()

		handler.Login(w, req, nil)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestLogout(t *testing.T) {
	mockSvc := new(MockAuthUsecase)
	cookieCfg := cookie.CookieConfig{
		RefreshName: "refresh_token",
	}

	handler := httpx.NewAuthHandler(mockSvc, cookieCfg, time.Hour, time.Hour*24, nil, nil)

	t.Run("Success Logout", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "valid_refresh"})
		w := httptest.NewRecorder()

		mockSvc.On("Logout", mock.Anything, "valid_refresh", mock.Anything).Return(nil).Once()

		handler.Logout(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminRevoke(t *testing.T) {
	mockSvc := new(MockAuthUsecase)
	handler := httpx.NewAuthHandler(mockSvc, cookie.CookieConfig{}, time.Hour, time.Hour*24, nil, nil)

	t.Run("Success Revoke", func(t *testing.T) {
		reqBody := httpx.AdminRevokeRequest{SessionId: "session123"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/revoke", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSvc.On("AdminRevokingSession", mock.Anything, "session123").Return(nil).Once()

		handler.AdminRevokeUser(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		reqBody := httpx.AdminRevokeRequest{SessionId: "missing"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/revoke", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSvc.On("AdminRevokingSession", mock.Anything, "missing").Return(coreerror.ErrNotFound).Once()

		handler.AdminRevokeUser(w, req, nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAuthMe(t *testing.T) {
	handler := httpx.NewAuthHandler(nil, cookie.CookieConfig{}, time.Hour, time.Hour*24, nil, nil)

	t.Run("Success AuthMe", func(t *testing.T) {
		actor := user.Actor{
			IdPengguna: 1,
			Username:   "testuser",
			Role:       user.ADMIN,
		}
		ctx := middleware.WithActor(context.Background(), actor)
		req := httptest.NewRequest(http.MethodGet, "/me", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		handler.AuthMe(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "testuser", data["username"])
	})

	t.Run("No Actor", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		w := httptest.NewRecorder()

		handler.AuthMe(w, req, nil)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
