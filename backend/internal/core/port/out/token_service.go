package out

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type AccessTokenService interface {
	GenerateAccessToken(idPengguna user.ID, role user.Role, username string, tokenDuration time.Duration) (string,error)
	VerifyAccessToken(token string, now time.Time) (idPengguna user.ID, role user.Role, username string,err error)
}

type RefreshTokenService interface {
	GenerateRefreshToken(sessionId string, tokenDuration time.Duration) (string,error)
	VerifyRefreshToken(token string, now time.Time) (sessionID string,err error)
}