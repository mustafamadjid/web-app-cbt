package integration_test

import (
	"testing"

	resetrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/reset_password"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/security/bcrypt"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	usersvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/reset_password"
	userrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetPasswordService_Reset(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := resetrepo.NewResetPasswordRepo(tx, nil)
	userRepo := userrepo.NewUserRepo(tx, nil)
	hasher := bcrypt.NewHasher(4)
	svc := usersvc.NewResetPasswordService(repo, hasher)

	email := user.Email("reset-password@example.com")
	id, err := userRepo.CreateUser(ctx, user.Pengguna{
		Username:       "reset_password_user",
		Email:          &email,
		PasswordHashed: "old-hash",
		NamaLengkap:    "Reset Password",
		JenisKelamin:   "LAKI_LAKI",
		Role:           user.GURU,
		StatusAkun:     user.AKTIF,
	})
	require.NoError(t, err)

	require.NoError(t, svc.ResetPasswordService(ctx, id, "new-password"))

	stored, err := userRepo.FindUserByID(ctx, id)
	require.NoError(t, err)
	assert.NotEqual(t, "old-hash", stored.PasswordHashed)
	assert.True(t, hasher.ComparePaswordAndHashed(stored.PasswordHashed, "new-password"))
}

