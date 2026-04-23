package userrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewUserRepo(tx, nil)

	email := user.Email("repo-user@example.com")
	noHP := "081234567890"
	id, err := repo.CreateUser(ctx, user.Pengguna{
		Username:       "repo_user_1",
		Email:          &email,
		PasswordHashed: "hash-1",
		NamaLengkap:    "Repo User",
		JenisKelamin:   "LAKI_LAKI",
		NoHp:           &noHP,
		Role:           user.GURU,
		StatusAkun:     user.AKTIF,
	})
	require.NoError(t, err)

	exists, err := repo.UserExistByUsername(ctx, "repo_user_1")
	require.NoError(t, err)
	assert.True(t, exists)

	found, err := repo.FindUserByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "repo_user_1", found.Username)
	assert.Equal(t, user.GURU, found.Role)

	updatedName := "Updated Repo User"
	updatedRole := user.SISWA
	updatedGender := "PEREMPUAN"
	err = repo.UpdateUser(ctx, id, updatepatch.Pengguna{
		NamaLengkap:  &updatedName,
		Role:         &updatedRole,
		JenisKelamin: &updatedGender,
	})
	require.NoError(t, err)

	updated, err := repo.FindUserByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.NamaLengkap)
	assert.Equal(t, updatedRole, updated.Role)
	assert.Equal(t, "PEREMPUAN", updated.JenisKelamin)

	items, err := repo.ListUser(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, items)

	err = repo.DeleteUser(ctx, id)
	require.NoError(t, err)

	_, err = repo.FindUserByID(ctx, id)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

func TestUserRepo_DeleteUsers(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewUserRepo(tx, nil)

	email1 := user.Email("bulk-user-1@example.com")
	email2 := user.Email("bulk-user-2@example.com")
	id1, err := repo.CreateUser(ctx, user.Pengguna{
		Username:       "bulk_user_1",
		Email:          &email1,
		PasswordHashed: "hash-1",
		NamaLengkap:    "Bulk User 1",
		JenisKelamin:   "LAKI_LAKI",
		Role:           user.GURU,
		StatusAkun:     user.AKTIF,
	})
	require.NoError(t, err)
	id2, err := repo.CreateUser(ctx, user.Pengguna{
		Username:       "bulk_user_2",
		Email:          &email2,
		PasswordHashed: "hash-2",
		NamaLengkap:    "Bulk User 2",
		JenisKelamin:   "PEREMPUAN",
		Role:           user.SISWA,
		StatusAkun:     user.AKTIF,
	})
	require.NoError(t, err)

	rows, err := repo.DeleteUsers(ctx, []user.ID{id1, id2})
	require.NoError(t, err)
	assert.EqualValues(t, 2, rows)
}

func TestUserRepo_CreateUser_DuplicateUsername(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewUserRepo(tx, nil)

	created := fixtures.CreateUser(user.ADMIN)
	email := user.Email("duplicate-user@example.com")

	_, err := repo.CreateUser(ctx, user.Pengguna{
		Username:       created.Username,
		Email:          &email,
		PasswordHashed: "hash-duplicate",
		NamaLengkap:    "Duplicate User",
		JenisKelamin:   "LAKI_LAKI",
		Role:           user.ADMIN,
		StatusAkun:     user.AKTIF,
	})
	assert.ErrorIs(t, err, coreerror.ErrUsernameTaken)
}
