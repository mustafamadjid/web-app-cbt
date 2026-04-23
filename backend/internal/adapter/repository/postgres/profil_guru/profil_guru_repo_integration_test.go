package profilgururepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfilGuruRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewProfilgGuruRepo(tx, nil)

	createdUser := fixtures.CreateUser(user.GURU)
	profilID, err := repo.CreateProfilGuru(ctx, user.ProfilGuru{
		IdPengguna:  createdUser.ID,
		Nip:         user.NIP("198001010000000123"),
		Jabatan:     "Guru Matematika",
		BidangStudi: "Matematika",
	})
	require.NoError(t, err)
	assert.NotZero(t, profilID)

	exists, err := repo.ExistByNIP(ctx, user.NIP("198001010000000123"))
	require.NoError(t, err)
	assert.True(t, exists)

	item, err := repo.FindProfilGuruByID(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.Username, item.Username)
	assert.Equal(t, "198001010000000123", item.Nip)

	updatedJabatan := "Wali Kelas"
	err = repo.UpdateProfilGuru(ctx, createdUser.ID, updatepatch.ProfilGuru{
		Jabatan: &updatedJabatan,
	})
	require.NoError(t, err)

	updated, err := repo.FindProfilGuruByID(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, updatedJabatan, updated.Jabatan)

	items, err := repo.GetListGuru(ctx, query.ListGuruFilter{Search: createdUser.Username, Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, createdUser.ID, items[0].IdPengguna)
}

func TestProfilGuruRepo_FindProfilGuruByID_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewProfilgGuruRepo(tx, nil)

	_, err := repo.FindProfilGuruByID(ctx, 999999)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
