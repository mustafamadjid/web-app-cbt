package integration_test

import (
	"testing"

	ruangrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ruang_ujian"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangcreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/create"
	ruangdelete "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/delete"
	ruangget "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/get"
	ruangupdate "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/update"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuangUjianService_CreateGetUpdateDelete(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := ruangrepo.NewRuangUjianRepo(tx, nil)
	createSvc := ruangcreate.NewRuangUjianService(repo)
	getSvc := ruangget.NewGetRuangUjianService(repo)
	updateSvc := ruangupdate.NewUpdateRuangUjianService(repo)
	deleteSvc := ruangdelete.NewDeleteRuangUjianService(repo)

	require.NoError(t, createSvc.CreateRuangUjianService(ctx, ruangujian.RuangUjian{
		KodeRuang:   " rg-int-01 ",
		NamaRuangan: "  Ruang Integration ",
	}))

	items, err := getSvc.GetRuangUjian(ctx, query.ListRuangUjianFilter{Search: "Integration", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, items)

	var idRuangan int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_ruangan FROM ruang_ujian WHERE kode_ruang = 'RG-INT-01'`).Scan(&idRuangan))

	itemByID, err := getSvc.GetRuangUjianById(ctx, idRuangan)
	require.NoError(t, err)
	assert.Equal(t, "Ruang Integration", itemByID.NamaRuangan)

	itemByCode, err := getSvc.GetRuangUjianByKode(ctx, "RG-INT-01")
	require.NoError(t, err)
	assert.Equal(t, idRuangan, int(itemByCode.IdRuangan))

	updatedCode := "rg-int-02"
	updatedName := "  Ruang Baru  "
	require.NoError(t, updateSvc.UpdateRuangUjian(ctx, idRuangan, updatepatch.UpdateRuangUjianPatch{
		KodeRuang: &updatedCode,
		NamaRuang: &updatedName,
	}))

	updated, err := getSvc.GetRuangUjianById(ctx, idRuangan)
	require.NoError(t, err)
	assert.Equal(t, "RG-INT-02", updated.KodeRuang)
	assert.Equal(t, "Ruang Baru", updated.NamaRuangan)

	require.NoError(t, deleteSvc.DeleteRuangUjian(ctx, idRuangan))
	_, err = getSvc.GetRuangUjianById(ctx, idRuangan)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

