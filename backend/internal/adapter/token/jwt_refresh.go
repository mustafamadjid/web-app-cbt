package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type JWTRefreshTokenService struct {
	secret []byte
	issuer string
}

func NewJWTRefreshTokenService(secret string, issuer string) *JWTRefreshTokenService {
	return &JWTRefreshTokenService{
		secret: []byte(secret),
		issuer: issuer,
	}
}

type refreshClaims struct {
	SessionId string `json:"session_id"`
	jwt.RegisteredClaims
}

func (s *JWTRefreshTokenService)GenerateRefreshToken(sessionId string, tokenDuration time.Duration) (string,error){
	now := time.Now()
	claims := refreshClaims{
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: s.issuer,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return tok.SignedString(s.secret)
}

func (s *JWTRefreshTokenService) VerifyRefreshToken(token string, now time.Time) (string, error) {
	parsed, err := jwt.ParseWithClaims(token, &refreshClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	}, jwt.WithTimeFunc(func() time.Time { return now }))
	if err != nil {
		return "", err
	}
	claims, ok := parsed.Claims.(*refreshClaims)
	if !ok || !parsed.Valid {
		return "", errors.New("invalid token")
	}
	if claims.SessionId == "" {
		return "", errors.New("missing session id")
	}
	return claims.SessionId, nil
}