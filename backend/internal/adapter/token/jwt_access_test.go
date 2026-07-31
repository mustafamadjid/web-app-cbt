package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTAccessTokenService(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		issuer string
	}{
		{name: "Path 1 -> constructs service with configured secret and issuer", secret: "access-secret", issuer: "cbt"},
		{name: "Path 2 -> accepts empty configuration", secret: "", issuer: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewJWTAccessTokenService(tt.secret, tt.issuer)

			require.NotNil(t, service)
			assert.Equal(t, []byte(tt.secret), service.secret)
			assert.Equal(t, tt.issuer, service.issuer)
		})
	}
}

func TestJWTAccessTokenService_GenerateAccessToken(t *testing.T) {
	service := NewJWTAccessTokenService("access-secret", "cbt")
	tests := []struct {
		name        string
		id          user.ID
		role        user.Role
		username    string
		duration    time.Duration
		wantExpired bool
	}{
		{name: "Path 1 -> generates token containing identity claims", id: 42, role: user.ADMIN, username: "admin", duration: time.Hour},
		{name: "Path 2 -> generates already expired token for negative duration", id: 7, role: user.SISWA, username: "siswa", duration: -time.Hour, wantExpired: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now().Add(-time.Second)
			signed, err := service.GenerateAccessToken(tt.id, tt.role, tt.username, tt.duration)
			after := time.Now().Add(time.Second)

			require.NoError(t, err)
			require.NotEmpty(t, signed)

			claims := &accesClaims{}
			_, err = jwt.ParseWithClaims(signed, claims, func(token *jwt.Token) (any, error) {
				return service.secret, nil
			}, jwt.WithoutClaimsValidation())
			require.NoError(t, err)
			assert.Equal(t, tt.id, claims.IdPengguna)
			assert.Equal(t, tt.role, claims.Role)
			assert.Equal(t, tt.username, claims.Username)
			assert.Equal(t, "cbt", claims.Issuer)
			assert.WithinRange(t, claims.IssuedAt.Time, before, after)
			assert.WithinRange(t, claims.NotBefore.Time, before, after)
			if tt.wantExpired {
				assert.True(t, claims.ExpiresAt.Before(claims.IssuedAt.Time))
			} else {
				assert.WithinDuration(t, claims.IssuedAt.Add(tt.duration), claims.ExpiresAt.Time, time.Second)
			}
		})
	}
}

func TestJWTAccessTokenService_VerifyAccessToken(t *testing.T) {
	service := NewJWTAccessTokenService("access-secret", "cbt")
	validToken, err := service.GenerateAccessToken(42, user.ADMIN, "admin", time.Hour)
	require.NoError(t, err)
	expiredToken, err := service.GenerateAccessToken(42, user.ADMIN, "admin", -time.Hour)
	require.NoError(t, err)
	wrongSecretToken, err := NewJWTAccessTokenService("other-secret", "cbt").GenerateAccessToken(42, user.ADMIN, "admin", time.Hour)
	require.NoError(t, err)
	wrongMethodToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, accesClaims{
		IdPengguna: 42, Role: user.ADMIN, Username: "admin",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
	}).SignedString([]byte("access-secret"))
	require.NoError(t, err)

	tests := []struct {
		name         string
		token        string
		wantID       user.ID
		wantRole     user.Role
		wantUsername string
		wantErr      bool
		now          time.Time
	}{
		{name: "Path 1 -> malformed token returns parse error", token: "not-a-jwt", wantErr: true},
		{name: "Path 2 -> non HS256 method is rejected", token: wrongMethodToken, wantErr: true},
		{name: "Path 3 -> invalid signature returns verification error", token: wrongSecretToken, wantErr: true},
		{name: "Path 4 -> expired token returns validation error", token: expiredToken, wantErr: true},
		{name: "Path 5 -> injected future time expires otherwise valid token", token: validToken, now: time.Now().Add(2 * time.Hour), wantErr: true},
		{name: "Path 6 -> valid token returns all identity claims", token: validToken, wantID: 42, wantRole: user.ADMIN, wantUsername: "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := tt.now
			if now.IsZero() {
				now = time.Now()
			}
			id, role, username, err := service.VerifyAccessToken(tt.token, now)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Zero(t, id)
				assert.Empty(t, role)
				assert.Empty(t, username)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, id)
			assert.Equal(t, tt.wantRole, role)
			assert.Equal(t, tt.wantUsername, username)
		})
	}
}
