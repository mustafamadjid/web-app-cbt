package ruangujianrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuangUjianRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewRuangUjianRepo(tx, nil)

	err := repo.CreateRuangUjian(ctx, ruangujian.RuangUjian{
		NamaRuangan: "Ruang Integration",
		KodeRuang:   "RG-IT-01",
	})
	require.NoError(t, err)

	var idRuangan int
	err = tx.QueryRow(ctx, `SELECT id_ruangan FROM ruang_ujian WHERE kode_ruang = 'RG-IT-01'`).Scan(&idRuangan)
	require.NoError(t, err)

	exists, err := repo.ExistByKodeRuang(ctx, "RG-IT-01")
	require.NoError(t, err)
	assert.True(t, exists)

	itemByID, err := repo.GetRuangUjianById(ctx, idRuangan)
	require.NoError(t, err)
	assert.Equal(t, "Ruang Integration", itemByID.NamaRuangan)

	itemByCode, err := repo.GetRuangUjianByKode(ctx, "RG-IT-01")
	require.NoError(t, err)
	assert.Equal(t, ruangujian.ID(idRuangan), itemByCode.IdRuangan)

	items, err := repo.GetRuangUjian(ctx, query.ListRuangUjianFilter{Search: "Integration", Limit: 10})
	require.NoError(t, err)
	assert.Len(t, items, 1)

	updatedName := "Ruang Updated"
	updatedCode := "RG-IT-02"
	err = repo.UpdateRuangUjian(ctx, idRuangan, updatepatch.UpdateRuangUjianPatch{
		KodeRuang: &updatedCode,
		NamaRuang: &updatedName,
	})
	require.NoError(t, err)

	updated, err := repo.GetRuangUjianById(ctx, idRuangan)
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.NamaRuangan)
	assert.Equal(t, updatedCode, updated.KodeRuang)

	err = repo.DeleteRuangUjian(ctx, idRuangan)
	require.NoError(t, err)

	_, err = repo.GetRuangUjianById(ctx, idRuangan)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
