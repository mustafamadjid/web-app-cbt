package out

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type TokenService interface {
	GenerateAccessToken(userID user.ID, role user.Role, tokenDuration time.Duration) (string,error)
	GenerateRefreshToken(sessionID string, tokenDuration time.Duration) (string,error)
	VerifyRefreshToken(token string) (sessionID string,err error)
}