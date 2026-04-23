package resetpasswordrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetPasswordRepo_ResetPassword(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewResetPasswordRepo(tx, nil)

	created := fixtures.CreateUser(user.ADMIN)

	err := repo.ResetPassword(ctx, created.ID, "new-hash-password")
	require.NoError(t, err)

	var stored string
	err = tx.QueryRow(ctx, `SELECT password FROM pengguna WHERE id_pengguna = $1`, created.ID).Scan(&stored)
	require.NoError(t, err)
	assert.Equal(t, "new-hash-password", stored)
}

func TestResetPasswordRepo_ResetPassword_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewResetPasswordRepo(tx, nil)

	err := repo.ResetPassword(ctx, user.ID(999999), "new-hash-password")
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
