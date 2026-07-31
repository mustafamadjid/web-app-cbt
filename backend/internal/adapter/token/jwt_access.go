package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type JWTAccessTokenService struct {
	secret []byte
	issuer string
}

func NewJWTAccessTokenService(secret string, issuer string) *JWTAccessTokenService {
	return &JWTAccessTokenService{
		secret: []byte(secret),
		issuer: issuer,
	}
}

type accesClaims struct {
	IdPengguna user.ID   `json:"id_pengguna"`
	Role       user.Role `json:"role"`
	Username   string    `json:"username"`
	jwt.RegisteredClaims
}

func (s *JWTAccessTokenService) GenerateAccessToken(idPengguna user.ID, role user.Role, username string, tokenDuration time.Duration) (string, error) {
	now := time.Now()
	claims := accesClaims{
		IdPengguna: idPengguna,
		Role:       role,
		Username:   username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.secret)
}

func (s *JWTAccessTokenService) VerifyAccessToken(token string, now time.Time) (idPengguna user.ID, role user.Role, username string, err error) {
	tok, err := jwt.ParseWithClaims(token, &accesClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	}, jwt.WithTimeFunc(func() time.Time { return now }))
	if err != nil {
		return 0, "", "", err
	}
	if claims, ok := tok.Claims.(*accesClaims); ok && tok.Valid {
		return claims.IdPengguna, claims.Role, claims.Username, nil
	}
	return 0, "", "", errors.New("invalid token")
}
