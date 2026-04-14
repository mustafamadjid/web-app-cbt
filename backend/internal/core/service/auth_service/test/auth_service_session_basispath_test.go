package auth_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	auth_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
	"github.com/stretchr/testify/assert"
)

func TestAuthServiceLogoutBasisPath(t *testing.T) {
	t.Parallel()

	baseStore := map[string]session.Session{
		"session_user_1": {
			SessionID: "session_user_1",
			UserID:    1,
			Role:      user.ADMIN,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	}

	tests := []struct {
		name      string
		token     string
		session   *FakeSessionRepo
		refresh   *FakeRefreshToken
		wantErr   error
		assertion func(*testing.T, *FakeSessionRepo)
	}{
		{
			name:    "path 1 -> refresh token tidak valid",
			token:   "invalid",
			session: &FakeSessionRepo{Store: map[string]session.Session{}},
			refresh: &FakeRefreshToken{},
			wantErr: coreerror.ErrInvalidToken,
			assertion: func(t *testing.T, repo *FakeSessionRepo) {
				assert.False(t, repo.Store["session_user_1"].Revoked)
			},
		},
		{
			name:    "path 2 -> revoke session gagal",
			token:   "REFRESH TOKEN : session_user_1",
			session: &FakeSessionRepo{Store: copySessionStore(baseStore), RevokeSessionErr: errors.New("revoke error")},
			refresh: &FakeRefreshToken{},
			wantErr: errors.New("revoke error"),
		},
		{
			name:    "path 3 -> logout berhasil",
			token:   "REFRESH TOKEN : session_user_1",
			session: &FakeSessionRepo{Store: copySessionStore(baseStore)},
			refresh: &FakeRefreshToken{},
			assertion: func(t *testing.T, repo *FakeSessionRepo) {
				assert.True(t, repo.Store["session_user_1"].Revoked)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			service := auth_service.NewAuthService(&FakeAuthUserRepo{}, &FakeUserRepo{}, &FakeHasher{}, tc.session, &FakeAccessToken{}, tc.refresh, 14*24*time.Hour)
			err := service.Logout(context.Background(), tc.token, time.Now())

			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			if tc.assertion != nil {
				tc.assertion(t, tc.session)
			}
		})
	}
}

func TestAuthServiceRefreshAccessTokenBasisPath(t *testing.T) {
	t.Parallel()

	baseUser := user.Pengguna{ID: 1, Username: "admin", Role: user.ADMIN}
	baseSession := session.Session{
		SessionID: "session_user_1",
		UserID:    1,
		Role:      user.ADMIN,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	tests := []struct {
		name      string
		token     string
		userRepo  *FakeUserRepo
		session   *FakeSessionRepo
		access    *FakeAccessToken
		refresh   *FakeRefreshToken
		wantErr   error
		wantToken string
	}{
		{
			name:     "path 1 -> token kosong",
			token:    "",
			userRepo: &FakeUserRepo{},
			session:  &FakeSessionRepo{},
			access:   &FakeAccessToken{},
			refresh:  &FakeRefreshToken{},
			wantErr:  coreerror.ErrNoTokenProvided,
		},
		{
			name:     "path 2 -> refresh token tidak valid",
			token:    "invalid",
			userRepo: &FakeUserRepo{},
			session:  &FakeSessionRepo{},
			access:   &FakeAccessToken{},
			refresh:  &FakeRefreshToken{},
			wantErr:  coreerror.ErrInvalidToken,
		},
		{
			name:     "path 3 -> get session gagal",
			token:    "REFRESH TOKEN : session_user_1",
			userRepo: &FakeUserRepo{},
			session:  &FakeSessionRepo{Store: map[string]session.Session{}},
			access:   &FakeAccessToken{},
			refresh:  &FakeRefreshToken{},
			wantErr:  coreerror.ErrNotFound,
		},
		{
			name:  "path 4 -> session expired",
			token: "REFRESH TOKEN : session_user_1",
			userRepo: &FakeUserRepo{
				ByID: map[user.ID]user.Pengguna{1: baseUser},
			},
			session: &FakeSessionRepo{
				Store: map[string]session.Session{
					"session_user_1": func() session.Session {
						s := baseSession
						s.ExpiresAt = time.Now().Add(-1 * time.Hour)
						return s
					}(),
				},
			},
			access:  &FakeAccessToken{},
			refresh: &FakeRefreshToken{},
			wantErr: coreerror.ErrSessionExpired,
		},
		{
			name:  "path 5 -> session revoked",
			token: "REFRESH TOKEN : session_user_1",
			userRepo: &FakeUserRepo{
				ByID: map[user.ID]user.Pengguna{1: baseUser},
			},
			session: &FakeSessionRepo{
				Store: map[string]session.Session{
					"session_user_1": func() session.Session {
						s := baseSession
						s.Revoked = true
						return s
					}(),
				},
			},
			access:  &FakeAccessToken{},
			refresh: &FakeRefreshToken{},
			wantErr: coreerror.ErrSessionExpired,
		},
		{
			name:  "path 6 -> find user gagal",
			token: "REFRESH TOKEN : session_user_1",
			userRepo: &FakeUserRepo{
				FindUserErr: errors.New("user error"),
			},
			session: &FakeSessionRepo{
				Store: map[string]session.Session{"session_user_1": baseSession},
			},
			access:  &FakeAccessToken{},
			refresh: &FakeRefreshToken{},
			wantErr: errors.New("user error"),
		},
		{
			name:  "path 7 -> generate access token gagal",
			token: "REFRESH TOKEN : session_user_1",
			userRepo: &FakeUserRepo{
				ByID: map[user.ID]user.Pengguna{1: baseUser},
			},
			session: &FakeSessionRepo{
				Store: map[string]session.Session{"session_user_1": baseSession},
			},
			access:  &FakeAccessToken{GenerateAccessTokenErr: errors.New("access error")},
			refresh: &FakeRefreshToken{},
			wantErr: errors.New("access error"),
		},
		{
			name:  "path 8 -> refresh access token berhasil",
			token: "REFRESH TOKEN : session_user_1",
			userRepo: &FakeUserRepo{
				ByID: map[user.ID]user.Pengguna{1: baseUser},
			},
			session: &FakeSessionRepo{
				Store: map[string]session.Session{"session_user_1": baseSession},
			},
			access:    &FakeAccessToken{},
			refresh:   &FakeRefreshToken{},
			wantToken: "ACCESS TOKEN :|1|ADMIN|admin",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			service := auth_service.NewAuthService(&FakeAuthUserRepo{}, tc.userRepo, &FakeHasher{}, tc.session, tc.access, tc.refresh, 14*24*time.Hour)
			got, err := service.RefreshAccessToken(context.Background(), tc.token, 15*time.Minute)

			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantToken, got)
		})
	}
}

func TestAuthServiceAdminRevokingSessionBasisPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		sessionID string
		session   *FakeSessionRepo
		wantErr   error
		assertion func(*testing.T, *FakeSessionRepo)
	}{
		{
			name:      "path 1 -> session id kosong",
			sessionID: "",
			session:   &FakeSessionRepo{Store: map[string]session.Session{}},
			wantErr:   coreerror.ErrNoSessionId,
		},
		{
			name:      "path 2 -> revoke session gagal",
			sessionID: "session_user_1",
			session: &FakeSessionRepo{
				Store:            map[string]session.Session{"session_user_1": {SessionID: "session_user_1", UserID: 1, Role: user.ADMIN, ExpiresAt: time.Now().Add(1 * time.Hour)}},
				RevokeSessionErr: errors.New("revoke error"),
			},
			wantErr: errors.New("revoke error"),
		},
		{
			name:      "path 3 -> admin revoke session berhasil",
			sessionID: "session_user_1",
			session: &FakeSessionRepo{
				Store: map[string]session.Session{"session_user_1": {SessionID: "session_user_1", UserID: 1, Role: user.ADMIN, ExpiresAt: time.Now().Add(1 * time.Hour)}},
			},
			assertion: func(t *testing.T, repo *FakeSessionRepo) {
				assert.True(t, repo.Store["session_user_1"].Revoked)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			service := auth_service.NewAuthService(&FakeAuthUserRepo{}, &FakeUserRepo{}, &FakeHasher{}, tc.session, &FakeAccessToken{}, &FakeRefreshToken{}, 14*24*time.Hour)
			err := service.AdminRevokingSession(context.Background(), tc.sessionID)

			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			if tc.assertion != nil {
				tc.assertion(t, tc.session)
			}
		})
	}
}

func copySessionStore(in map[string]session.Session) map[string]session.Session {
	out := make(map[string]session.Session, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
