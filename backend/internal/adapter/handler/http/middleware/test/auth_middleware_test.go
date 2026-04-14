package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	middleware "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithActorAndActorFromContext(t *testing.T) {
	actor := user.Actor{
		IdPengguna: 12,
		Role:       user.GURU,
		Username:   "guru1",
	}

	ctx := middleware.WithActor(context.Background(), actor)
	got, ok := middleware.ActorFromContext(ctx)

	require.True(t, ok)
	assert.Equal(t, actor, got)
}

func TestRequireValidTokenAndSession_Success(t *testing.T) {
	access := &fakeAccessTokenService{
		userID:   12,
		role:     user.ADMIN,
		username: "admin",
	}
	sessionRepo := &fakeSessionRepository{hasActive: true}
	cfg := cookie.CookieConfig{AccessName: "access_token", RefreshName: "refresh_token"}

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		actor, ok := middleware.ActorFromContext(r.Context())
		require.True(t, ok)
		assert.Equal(t, user.ID(12), actor.IdPengguna)
		assert.Equal(t, user.ADMIN, actor.Role)
		assert.Equal(t, "admin", actor.Username)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cfg.AccessName, Value: "valid-token"})
	rec := httptest.NewRecorder()

	handler := middleware.RequireValidTokenAndSession(next, access, &fakeRefreshTokenService{}, sessionRepo, cfg)
	handler.ServeHTTP(rec, req)

	assert.True(t, nextCalled)
	assert.Equal(t, "valid-token", access.token)
	assert.Equal(t, user.ID(12), sessionRepo.userID)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireValidTokenAndSession_MissingCookie(t *testing.T) {
	access := &fakeAccessTokenService{}
	sessionRepo := &fakeSessionRepository{hasActive: true}
	cfg := cookie.CookieConfig{AccessName: "access_token", RefreshName: "refresh_token"}

	nextCalled := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler := middleware.RequireValidTokenAndSession(next, access, &fakeRefreshTokenService{}, sessionRepo, cfg)
	handler.ServeHTTP(rec, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRequireValidTokenAndSession_InvalidTokenClearsCookies(t *testing.T) {
	access := &fakeAccessTokenService{err: assert.AnError}
	sessionRepo := &fakeSessionRepository{hasActive: true}
	cfg := cookie.CookieConfig{AccessName: "access_token", RefreshName: "refresh_token"}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cfg.AccessName, Value: "invalid-token"})
	rec := httptest.NewRecorder()

	handler := middleware.RequireValidTokenAndSession(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), access, &fakeRefreshTokenService{}, sessionRepo, cfg)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 2)
	assert.ElementsMatch(t, []string{cfg.AccessName, cfg.RefreshName}, []string{cookies[0].Name, cookies[1].Name})
}

func TestRequireValidTokenAndSession_InactiveSessionClearsCookies(t *testing.T) {
	access := &fakeAccessTokenService{
		userID:   99,
		role:     user.GURU,
		username: "guru99",
	}
	sessionRepo := &fakeSessionRepository{hasActive: false}
	cfg := cookie.CookieConfig{AccessName: "access_token", RefreshName: "refresh_token"}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cfg.AccessName, Value: "valid-token"})
	rec := httptest.NewRecorder()

	handler := middleware.RequireValidTokenAndSession(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), access, &fakeRefreshTokenService{}, sessionRepo, cfg)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 2)
	assert.ElementsMatch(t, []string{cfg.AccessName, cfg.RefreshName}, []string{cookies[0].Name, cookies[1].Name})
}

func TestRequireActorRole(t *testing.T) {
	t.Run("allowed", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{
			IdPengguna: 1,
			Role:       user.ADMIN,
			Username:   "admin",
		}))
		rec := httptest.NewRecorder()

		middleware.RequireActorRole(next, user.ADMIN, user.GURU).ServeHTTP(rec, req)

		assert.True(t, nextCalled)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{
			IdPengguna: 1,
			Role:       user.SISWA,
			Username:   "siswa",
		}))
		rec := httptest.NewRecorder()

		middleware.RequireActorRole(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), user.ADMIN, user.GURU).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("missing actor", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		middleware.RequireActorRole(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), user.ADMIN).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
