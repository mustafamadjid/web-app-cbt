package matapelajaranrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapelRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewMapelRepo(tx, nil)

	kelas := fixtures.CreateKelas(66)
	err := repo.CreateMapel(ctx, matapelajaran.MataPelajaran{
		IdKelas:   matapelajaran.ID(kelas.ID),
		KodeMapel: "KODE66",
		NamaMapel: "Matapelajaran IT",
		Deskripsi: "Deskripsi mapel",
	})
	require.NoError(t, err)

	var idMapel int
	err = tx.QueryRow(ctx, `SELECT id_mapel FROM mata_pelajaran WHERE kode_mapel = 'KODE66'`).Scan(&idMapel)
	require.NoError(t, err)

	exists, err := repo.ExistKodeMapel(ctx, "KODE66")
	require.NoError(t, err)
	assert.True(t, exists)

	item, err := repo.GetMapelById(ctx, idMapel)
	require.NoError(t, err)
	assert.Equal(t, "Matapelajaran IT", item.NamaMapel)

	found, err := repo.GetMapel(ctx, query.ListMapelFilter{Search: "IT", Limit: 10})
	require.NoError(t, err)
	assert.Len(t, found, 1)

	updatedName := "Mapel Updated"
	updatedDesc := "Desc Updated"
	err = repo.UpdateMapel(ctx, idMapel, updatepatch.UpdateMapelPatch{
		NamaMapel: &updatedName,
		Deskripsi: &updatedDesc,
	})
	require.NoError(t, err)

	updated, err := repo.GetMapelById(ctx, idMapel)
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.NamaMapel)
	assert.Equal(t, updatedDesc, updated.Deskripsi)

	err = repo.DeleteMapel(ctx, idMapel)
	require.NoError(t, err)

	_, err = repo.GetMapelById(ctx, idMapel)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
