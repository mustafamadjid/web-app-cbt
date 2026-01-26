package auth_service

import (
	"context"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/login"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth_port_out"
)

type AuthService struct {
	authUsers auth_port_out.AuthUserrepository
	hasher out.PasswordHasher
	sessions out.SessionRepository
	accessTokens out.AccessTokenService
	refreshTokens out.RefreshTokenService
}

func NewAuthService(authUser auth_port_out.AuthUserrepository, hash out.PasswordHasher, session out.SessionRepository, accessToken out.AccessTokenService, refreshToken out.RefreshTokenService) *AuthService {
	return &AuthService{authUsers: authUser, hasher: hash, sessions: session, accessTokens: accessToken, refreshTokens: refreshToken}
}


func (authService *AuthService) Login(ctx context.Context, cmd login.LoginCmd) ( login.LoginRes, error) {
	u,err := authService.authUsers.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return login.LoginRes{}, coreerror.ErrInvalidCreds
	}

	if u.StatusAkun != user.AKTIF {
		return  login.LoginRes{}, coreerror.ErrInvalidCreds
	}

	if !authService.hasher.ComparePaswordAndHashed(u.PasswordHashed, cmd.Password) {
		return  login.LoginRes{}, coreerror.ErrInvalidCreds
	}

	refreshExp := time.Now().Add(14 * 24 * time.Hour)
	sessionId, err := authService.sessions.CreateSession(ctx, u.ID, refreshExp)
	if err != nil {
		return  login.LoginRes{}, err
	}

	accessToken, err := authService.accessTokens.GenerateAccessToken(u.ID, u.Role,u.Username, time.Minute * 15)
	if err != nil {
		return  login.LoginRes{}, err
	}

	refreshToken, err := authService.refreshTokens.GenerateRefreshToken(sessionId, time.Hour * 24 * 14)
	if err != nil {
		return  login.LoginRes{}, err
	}

	return  login.LoginRes{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (authService *AuthService) Logout(ctx context.Context, refreshtoken string, now time.Time) error {
	sid,err := authService.refreshTokens.VerifyRefreshToken(refreshtoken,time.Time{})
	if err != nil {
		return coreerror.ErrInvalidToken
	}

	if err :=authService.sessions.RevokeSession(context.Background(),sid); err != nil {
		return err
	}
	return nil
}
