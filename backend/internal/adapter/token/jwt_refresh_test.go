package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTRefreshTokenService(t *testing.T) {
	tests := []struct {
		name, secret, issuer string
	}{
		{name: "Path 1 -> constructs service with configured values", secret: "refresh-secret", issuer: "cbt"},
		{name: "Path 2 -> accepts empty configuration"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewJWTRefreshTokenService(tt.secret, tt.issuer)
			require.NotNil(t, service)
			assert.Equal(t, []byte(tt.secret), service.secret)
			assert.Equal(t, tt.issuer, service.issuer)
		})
	}
}

func TestJWTRefreshTokenService_GenerateRefreshToken(t *testing.T) {
	service := NewJWTRefreshTokenService("refresh-secret", "cbt")
	tests := []struct {
		name, sessionID string
		duration        time.Duration
	}{
		{name: "Path 1 -> generates token containing session claim", sessionID: "session-1", duration: 24 * time.Hour},
		{name: "Path 2 -> generates token with empty session for verifier validation", duration: time.Hour},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signed, err := service.GenerateRefreshToken(tt.sessionID, tt.duration)
			require.NoError(t, err)
			claims := &refreshClaims{}
			_, err = jwt.ParseWithClaims(signed, claims, func(token *jwt.Token) (any, error) {
				return service.secret, nil
			}, jwt.WithoutClaimsValidation())
			require.NoError(t, err)
			assert.Equal(t, tt.sessionID, claims.SessionId)
			assert.Equal(t, "cbt", claims.Issuer)
			assert.WithinDuration(t, claims.IssuedAt.Add(tt.duration), claims.ExpiresAt.Time, time.Second)
		})
	}
}

func TestJWTRefreshTokenService_VerifyRefreshToken(t *testing.T) {
	service := NewJWTRefreshTokenService("refresh-secret", "cbt")
	now := time.Now().UTC().Truncate(time.Second)
	sign := func(method jwt.SigningMethod, claims refreshClaims, secret string) string {
		t.Helper()
		token, err := jwt.NewWithClaims(method, claims).SignedString([]byte(secret))
		require.NoError(t, err)
		return token
	}
	baseClaims := func(sessionID string, expiresAt time.Time) refreshClaims {
		return refreshClaims{SessionId: sessionID, RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "cbt", IssuedAt: jwt.NewNumericDate(now.Add(-time.Minute)),
			NotBefore: jwt.NewNumericDate(now.Add(-time.Minute)), ExpiresAt: jwt.NewNumericDate(expiresAt),
		}}
	}

	tests := []struct {
		name, token, wantSession string
		wantErr                  string
	}{
		{name: "Path 1 -> malformed token returns parse error", token: "invalid", wantErr: "token is malformed"},
		{name: "Path 2 -> non HS256 method is rejected", token: sign(jwt.SigningMethodHS384, baseClaims("session-1", now.Add(time.Hour)), "refresh-secret"), wantErr: "unexpected signing method"},
		{name: "Path 3 -> wrong signature returns verification error", token: sign(jwt.SigningMethodHS256, baseClaims("session-1", now.Add(time.Hour)), "wrong-secret"), wantErr: "signature is invalid"},
		{name: "Path 4 -> expired token at injected time returns validation error", token: sign(jwt.SigningMethodHS256, baseClaims("session-1", now.Add(-time.Second)), "refresh-secret"), wantErr: "token is expired"},
		{name: "Path 5 -> valid token without session id is rejected", token: sign(jwt.SigningMethodHS256, baseClaims("", now.Add(time.Hour)), "refresh-secret"), wantErr: "missing session id"},
		{name: "Path 6 -> valid token returns session id", token: sign(jwt.SigningMethodHS256, baseClaims("session-1", now.Add(time.Hour)), "refresh-secret"), wantSession: "session-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionID, err := service.VerifyRefreshToken(tt.token, now)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
				assert.Empty(t, sessionID)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantSession, sessionID)
		})
	}
}
