package auth

import (
	"context"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	out_auth "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth"
)

type AuthService struct {
	authUsers out_auth.AuthUserrepository
	hasher out.PasswordHasher
	sessions out.SessionRepository
	tokens out.TokenService
}

func NewAuthService(authUser out_auth.AuthUserrepository, hash out.PasswordHasher, session out.SessionRepository, token out.TokenService) *AuthService {
	return &AuthService{authUsers: authUser, hasher: hash, sessions: session, tokens: token}
}
type LoginCmd struct {
	Username    string
	Password string
}

type LoginRes struct {
	AccessToken  string
	RefreshToken string
}

func (authService *AuthService) Login(ctx context.Context, cmd LoginCmd) (LoginRes, error) {
	user,err := authService.authUsers.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return LoginRes{}, coreerror.ErrInvalidCreds
	}

	if user.StatusAkun != "AKTIF" {
		return LoginRes{}, coreerror.ErrInvalidCreds
	}

	if !authService.hasher.ComparePaswordAndHashed(user.PasswordHashed, cmd.Password) {
		return LoginRes{}, coreerror.ErrInvalidCreds
	}

	refreshExp := time.Now().Add(14 * 24 * time.Hour)
	sessionId, err := authService.sessions.CreateSession(ctx, user.ID, refreshExp)
	if err != nil {
		return LoginRes{}, err
	}

	accessToken, err := authService.tokens.GenerateAccessToken(user.ID, user.Role, time.Minute * 15)
	if err != nil {
		return LoginRes{}, err
	}

	refreshToken, err := authService.tokens.GenerateRefreshToken(sessionId, time.Hour * 24 * 14)
	if err != nil {
		return LoginRes{}, err
	}

	return LoginRes{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}