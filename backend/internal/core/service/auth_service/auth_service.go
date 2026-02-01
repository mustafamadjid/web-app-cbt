package auth_service

import (
	"context"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"

	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth_port_out"
	
)

type AuthService struct {
	authUsers auth_port_out.AuthUserrepository
	users 		outuser.UserRepository
	hasher out.PasswordHasher
	sessions out.SessionRepository
	accessTokens out.AccessTokenService
	refreshTokens out.RefreshTokenService
}

func NewAuthService(authUser auth_port_out.AuthUserrepository,user outuser.UserRepository, hash out.PasswordHasher, session out.SessionRepository, accessToken out.AccessTokenService, refreshToken out.RefreshTokenService) *AuthService {
	return &AuthService{authUsers: authUser,users: user, hasher: hash, sessions: session, accessTokens: accessToken, refreshTokens: refreshToken}
}


func (authService *AuthService) Login(ctx context.Context, cmd LoginCmd) ( LoginRes, error) {
	u,err := authService.authUsers.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return LoginRes{}, coreerror.ErrInvalidCreds
	}

	if u.StatusAkun != user.AKTIF {
		return  LoginRes{}, coreerror.ErrInvalidCreds
	}

	if !authService.hasher.ComparePaswordAndHashed(u.PasswordHashed, cmd.Password) {
		return  LoginRes{}, coreerror.ErrInvalidCreds
	}

	checkSession, err := authService.sessions.HasActiveSession(ctx,u.ID)
	if err != nil {
		return  LoginRes{}, err
	}

	if checkSession {
		return  LoginRes{}, coreerror.ErrHasSession
	}

	
	refreshExp := time.Now().Add(14 * 24 * time.Hour)
	sessionId, err := authService.sessions.CreateSession(ctx, u.ID, refreshExp)
	if err != nil {
		return  LoginRes{}, err
	}

	accessToken, err := authService.accessTokens.GenerateAccessToken(u.ID, u.Role,u.Username, time.Minute * 15)
	if err != nil {
		return  LoginRes{}, err
	}

	refreshToken, err := authService.refreshTokens.GenerateRefreshToken(sessionId, time.Hour * 24 * 14)
	if err != nil {
		return  LoginRes{}, err
	}

	return  LoginRes{IdPengguna: u.ID,Username:u.Username ,AccessToken: accessToken, RefreshToken: refreshToken}, nil
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

func(authService *AuthService)RefreshAccessToken(ctx context.Context, refreshToken string, accessTTL time.Duration) (string, error){
	 if refreshToken == "" {
		return "", coreerror.ErrNoTokenProvided
	 }

	 now := time.Now()
	sid,err := authService.refreshTokens.VerifyRefreshToken(refreshToken, now)
	if err != nil {
		return "", err
	}

	sess, err := authService.sessions.GetSession(ctx, sid)
	if err != nil {
		return "", err
	}

	if now.After(sess.ExpiresAt) || sess.Revoked {
		return "", coreerror.ErrSessionExpired 
	}

	u, err := authService.users.FindUserByID(ctx,sess.UserID)
	if err != nil {
		return "", err
	}

	accessToken, err := authService.accessTokens.GenerateAccessToken(sess.UserID,u.Role, u.Username, accessTTL)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func(authService *AuthService)AdminRevokingSession(ctx context.Context, sessionID string) error{
	if sessionID == "" {
		return coreerror.ErrNoSessionId
	}
	if err := authService.sessions.RevokeSession(ctx,sessionID); err != nil {
		return err
	}
	return nil
}

