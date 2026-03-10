package auth_service

import (
	"context"
	"strings"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"

	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth_port_out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type AuthService struct {
	authUsers     auth_port_out.AuthUserrepository
	users         outuser.UserRepository
	hasher        out.PasswordHasher
	sessions      out.SessionRepository
	accessTokens  out.AccessTokenService
	refreshTokens out.RefreshTokenService

	refreshTTL time.Duration
}

func NewAuthService(authUser auth_port_out.AuthUserrepository, user outuser.UserRepository, hash out.PasswordHasher, session out.SessionRepository, accessToken out.AccessTokenService, refreshToken out.RefreshTokenService, refreshTTL time.Duration) *AuthService {
	return &AuthService{authUsers: authUser, users: user, hasher: hash, sessions: session, accessTokens: accessToken, refreshTokens: refreshToken, refreshTTL: refreshTTL}
}

func (authService *AuthService) Login(ctx context.Context, cmd LoginCmd) (LoginRes, error) {
	logger := corelog.FromContext(ctx)
	cmd.Username = strings.TrimSpace(cmd.Username)
	if err := user.CheckUsernameLength(cmd.Username); err != nil {
		return LoginRes{}, err
	}

	u, err := authService.authUsers.FindByUsername(ctx, cmd.Username)
	if err != nil {
		logger.Info(ctx, "invalid credentials", "layer", "core.service", "op", "auth.login", "err", err)
		return LoginRes{}, coreerror.ErrInvalidCreds
	}

	if u.StatusAkun != user.AKTIF {
		return LoginRes{}, coreerror.ErrInvalidCreds
	}

	if !authService.hasher.ComparePaswordAndHashed(u.PasswordHashed, cmd.Password) {
		return LoginRes{}, coreerror.ErrInvalidCreds
	}

	_, err = authService.sessions.RevokeExpiredSessions(ctx, u.ID)
	if err != nil {
		logger.Error(ctx, "failed revoking expired sessions", "layer", "core.service", "op", "auth.login", "user_id", u.ID, "err", err)
		return LoginRes{}, err
	}

	checkSession, err := authService.sessions.HasActiveSession(ctx, u.ID)
	if err != nil {
		logger.Error(ctx, "failed checking active session", "layer", "core.service", "op", "auth.login", "user_id", u.ID, "err", err)
		return LoginRes{}, err
	}

	if checkSession {
		return LoginRes{}, coreerror.ErrHasSession
	}

	refreshExp := time.Now().Add(authService.refreshTTL)
	sessionId, err := authService.sessions.CreateSession(ctx, u.ID, refreshExp)
	if err != nil {
		logger.Error(ctx, "failed creating session", "layer", "core.service", "op", "auth.login", "user_id", u.ID, "err", err)
		return LoginRes{}, err
	}

	accessToken, err := authService.accessTokens.GenerateAccessToken(u.ID, u.Role, u.Username, time.Minute*15)
	if err != nil {
		logger.Error(ctx, "failed generating access token", "layer", "core.service", "op", "auth.login", "user_id", u.ID, "err", err)
		return LoginRes{}, err
	}

	refreshToken, err := authService.refreshTokens.GenerateRefreshToken(sessionId, time.Hour*24*14)
	if err != nil {
		logger.Error(ctx, "failed generating refresh token", "layer", "core.service", "op", "auth.login", "session_id", sessionId, "err", err)
		return LoginRes{}, err
	}

	return LoginRes{IdPengguna: u.ID, Username: u.Username, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (authService *AuthService) Logout(ctx context.Context, refreshtoken string, now time.Time) error {
	logger := corelog.FromContext(ctx)
	sid, err := authService.refreshTokens.VerifyRefreshToken(refreshtoken, now)
	if err != nil {
		logger.Info(ctx, "invalid refresh token", "layer", "core.service", "op", "auth.logout", "err", err)
		return coreerror.ErrInvalidToken
	}

	if err := authService.sessions.RevokeSession(context.Background(), sid); err != nil {
		logger.Error(ctx, "failed revoking session", "layer", "core.service", "op", "auth.logout", "session_id", sid, "err", err)
		return err
	}
	return nil
}

func (authService *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string, accessTTL time.Duration) (string, error) {
	logger := corelog.FromContext(ctx)
	if refreshToken == "" {
		return "", coreerror.ErrNoTokenProvided
	}

	now := time.Now()
	sid, err := authService.refreshTokens.VerifyRefreshToken(refreshToken, now)
	if err != nil {
		logger.Info(ctx, "invalid refresh token", "layer", "core.service", "op", "auth.refresh", "err", err)
		return "", err
	}

	sess, err := authService.sessions.GetSession(ctx, sid)
	if err != nil {
		logger.Error(ctx, "failed getting session", "layer", "core.service", "op", "auth.refresh", "session_id", sid, "err", err)
		return "", err
	}

	if now.After(sess.ExpiresAt) || sess.Revoked {
		return "", coreerror.ErrSessionExpired
	}

	u, err := authService.users.FindUserByID(ctx, sess.UserID)
	if err != nil {
		logger.Error(ctx, "failed finding user", "layer", "core.service", "op", "auth.refresh", "user_id", sess.UserID, "err", err)
		return "", err
	}

	accessToken, err := authService.accessTokens.GenerateAccessToken(sess.UserID, u.Role, u.Username, accessTTL)
	if err != nil {
		logger.Error(ctx, "failed generating access token", "layer", "core.service", "op", "auth.refresh", "user_id", sess.UserID, "err", err)
		return "", err
	}

	return accessToken, nil
}

func (authService *AuthService) AdminRevokingSession(ctx context.Context, sessionID string) error {
	logger := corelog.FromContext(ctx)
	if sessionID == "" {
		return coreerror.ErrNoSessionId
	}
	if err := authService.sessions.RevokeSession(ctx, sessionID); err != nil {
		logger.Error(ctx, "failed revoking session", "layer", "core.service", "op", "auth.admin_revoke", "session_id", sessionID, "err", err)
		return err
	}
	return nil
}
