package integration_test

import (
	"testing"
	"time"

	authuserrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/auth_user"
	sessionrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/session"
	userrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/user"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/security/bcrypt"
	jwttoken "github.com/mustafamadjid/web-app-cbt/internal/adapter/token"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	authsvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_LoginRefreshLogoutAndAdminRevoke(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	userRepo := userrepo.NewUserRepo(tx, nil)
	authRepo := authuserrepo.NewAuthUserRepo(tx, nil)
	sessionRepo := sessionrepo.NewSessionRepo(tx, nil)
	hasher := bcrypt.NewHasher(4)
	accessSvc := jwttoken.NewJWTAccessTokenService("secret-access", "integration")
	refreshSvc := jwttoken.NewJWTRefreshTokenService("secret-refresh", "integration")
	svc := authsvc.NewAuthService(authRepo, userRepo, hasher, sessionRepo, accessSvc, refreshSvc, 24*time.Hour)

	email := user.Email("auth-integration@example.com")
	passHash, err := hasher.GenerateHash("password-123")
	require.NoError(t, err)
	id, err := userRepo.CreateUser(ctx, user.Pengguna{
		Username:       "auth_user_01",
		Email:          &email,
		PasswordHashed: passHash,
		NamaLengkap:    "Auth Integration",
		JenisKelamin:   "LAKI_LAKI",
		Role:           user.ADMIN,
		StatusAkun:     user.AKTIF,
	})
	require.NoError(t, err)
	_ = id

	loginRes, err := svc.Login(ctx, authsvc.LoginCmd{Username: " auth_user_01 ", Password: "password-123"})
	require.NoError(t, err)
	assert.NotEmpty(t, loginRes.AccessToken)
	assert.NotEmpty(t, loginRes.RefreshToken)

	accessID, accessRole, accessUsername, err := accessSvc.VerifyAccessToken(loginRes.AccessToken, time.Now())
	require.NoError(t, err)
	assert.Equal(t, user.ADMIN, accessRole)
	assert.Equal(t, "auth_user_01", accessUsername)
	assert.Equal(t, id, accessID)

	newAccess, err := svc.RefreshAccessToken(ctx, loginRes.RefreshToken, time.Minute*30)
	require.NoError(t, err)
	assert.NotEmpty(t, newAccess)

	sid, err := refreshSvc.VerifyRefreshToken(loginRes.RefreshToken, time.Now())
	require.NoError(t, err)
	require.NoError(t, svc.AdminRevokingSession(ctx, sid))

	session, err := sessionRepo.GetSession(ctx, sid)
	require.NoError(t, err)
	assert.True(t, session.Revoked)

	loginRes2, err := svc.Login(ctx, authsvc.LoginCmd{Username: "auth_user_01", Password: "password-123"})
	require.NoError(t, err)
	require.NoError(t, svc.Logout(ctx, loginRes2.RefreshToken, time.Now()))
}
