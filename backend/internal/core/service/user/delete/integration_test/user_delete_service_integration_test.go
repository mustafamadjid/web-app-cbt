package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	deletefilesvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/delete_file_system"
	userrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/user"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	usersvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDeleteService_DeleteAndDeleteMany(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := userrepo.NewUserRepo(tx, nil)
	tmpDir := t.TempDir()
	deleteFileSvc := deletefilesvc.NewDeleteFileService(tmpDir)
	svc := usersvc.NewDeleteUserService(repo, deleteFileSvc)

	foto := "/uploads/user-delete.png"
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "user-delete.png"), []byte("photo"), 0o644))

	email := user.Email("delete-user@example.com")
	noHP := "081234567899"
	id1, err := repo.CreateUser(ctx, user.Pengguna{
		Username:       "delete_user_01",
		Email:          &email,
		PasswordHashed: "hashed-delete-1",
		NamaLengkap:    "Delete User 1",
		JenisKelamin:   "LAKI_LAKI",
		NoHp:           &noHP,
		Foto:           foto,
		Role:           user.GURU,
		StatusAkun:     user.AKTIF,
	})
	require.NoError(t, err)

	require.NoError(t, svc.Delete(ctx, id1))
	_, err = repo.FindUserByID(ctx, id1)
	assert.Error(t, err)
	_, err = os.Stat(filepath.Join(tmpDir, "user-delete.png"))
	assert.Error(t, err)

	email2 := user.Email("delete-many-1@example.com")
	email3 := user.Email("delete-many-2@example.com")
	id2, err := repo.CreateUser(ctx, user.Pengguna{Username: "delete_user_02", Email: &email2, PasswordHashed: "hashed-2", NamaLengkap: "Delete User 2", JenisKelamin: "LAKI_LAKI", Role: user.GURU, StatusAkun: user.AKTIF})
	require.NoError(t, err)
	id3, err := repo.CreateUser(ctx, user.Pengguna{Username: "delete_user_03", Email: &email3, PasswordHashed: "hashed-3", NamaLengkap: "Delete User 3", JenisKelamin: "PEREMPUAN", Role: user.SISWA, StatusAkun: user.AKTIF})
	require.NoError(t, err)

	affected, err := svc.DeleteMany(ctx, []user.ID{id2, id3})
	require.NoError(t, err)
	assert.EqualValues(t, 2, affected)
}

