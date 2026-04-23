package authuserrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthUserRepo_FindByUsername(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewAuthUserRepo(tx, nil)

	created := fixtures.CreateUser(user.GURU)

	found, err := repo.FindByUsername(ctx, created.Username)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.Username, found.Username)
	assert.Equal(t, user.GURU, found.Role)
	assert.Equal(t, user.AKTIF, found.StatusAkun)
	assert.NotEmpty(t, found.PasswordHashed)
}

func TestAuthUserRepo_FindByUsername_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewAuthUserRepo(tx, nil)

	_, err := repo.FindByUsername(ctx, "missing-user")
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
