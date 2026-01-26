package out

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type AccessTokenService interface {
	GenerateAccessToken(userID user.ID, role user.Role, username string, tokenDuration time.Duration) (string,error)
	VerifyAccessToken(token string, now time.Time) (userID user.ID, role user.Role, username string,err error)
}

type RefreshTokenService interface {
	GenerateRefreshToken(sessionID string, tokenDuration time.Duration) (string,error)
	VerifyRefreshToken(token string, now time.Time) (sessionID string,err error)
}