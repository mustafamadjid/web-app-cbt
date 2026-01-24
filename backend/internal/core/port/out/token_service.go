package out

import "time"

type TokenService interface {
	GenerateAccessToken(userID int, role string, tokenDuration time.Duration) (string,error)
	GenerateRefreshToken(sessionID string, tokenDuration time.Duration) (string,error)
	VerifyRefreshToken(token string) (sessionID string,err error)
}