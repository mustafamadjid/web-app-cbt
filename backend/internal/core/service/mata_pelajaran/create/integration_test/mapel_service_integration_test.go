package integration_test

import (
	"testing"

	mapelrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/mata_pelajaran"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	mapelcreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/create"
	mapeldelete "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/delete"
	mapelget "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/get"
	mapelupdate "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapelService_CreateGetUpdateDelete(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	kelas := fixtures.CreateKelas(8600)

	repo := mapelrepo.NewMapelRepo(tx, nil)
	createSvc := mapelcreate.NewMapelService(repo)
	getSvc := mapelget.NewGetMapelService(repo)
	updateSvc := mapelupdate.NewUpdateMapelService(repo)
	deleteSvc := mapeldelete.NewDeleteMapelService(repo)

	require.NoError(t, createSvc.CreateMapelService(ctx, matapelajaran.MataPelajaran{
		IdKelas:   matapelajaran.ID(kelas.ID),
		KodeMapel: " mp-int-01 ",
		NamaMapel: "  Mapel Integration ",
		Deskripsi: "  Deskripsi Mapel  ",
	}))

	items, err := getSvc.GetMapelService(ctx, query.ListMapelFilter{Search: "Integration", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, items)

	var idMapel int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_mapel FROM mata_pelajaran WHERE kode_mapel = 'MP-INT-01'`).Scan(&idMapel))

	found, err := getSvc.GetMapelById(ctx, idMapel)
	require.NoError(t, err)
	assert.Equal(t, "Mapel Integration", found.NamaMapel)

	newName := "Mapel Baru"
	newDesc := "Deskripsi Baru"
	require.NoError(t, updateSvc.UpdateMapelService(ctx, idMapel, updatepatch.UpdateMapelPatch{
		NamaMapel: &newName,
		Deskripsi: &newDesc,
	}))

	updated, err := getSvc.GetMapelById(ctx, idMapel)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.NamaMapel)
	assert.Equal(t, newDesc, updated.Deskripsi)

	require.NoError(t, deleteSvc.DeleteMapelService(ctx, idMapel))
	_, err = getSvc.GetMapelById(ctx, idMapel)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
