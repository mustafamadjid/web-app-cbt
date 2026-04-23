package sesirepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSesiRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewSesirepo(tx, nil)

	err := repo.CreateSesi(ctx, sesi.Sesi{
		KodeSesi: "S-INT-01",
		NamaSesi: "Sesi Integration",
	})
	require.NoError(t, err)

	var idSesi int
	err = tx.QueryRow(ctx, `SELECT id_sesi FROM sesi_ujian WHERE kode_sesi = 'S-INT-01'`).Scan(&idSesi)
	require.NoError(t, err)

	exists, err := repo.ExistByKodeSesi(ctx, "s-int-01")
	require.NoError(t, err)
	assert.True(t, exists)

	itemByID, err := repo.GetSesiById(ctx, idSesi)
	require.NoError(t, err)
	assert.Equal(t, "Sesi Integration", itemByID.NamaSesi)

	itemByCode, err := repo.GetSesiByKode(ctx, "S-INT-01")
	require.NoError(t, err)
	assert.Equal(t, sesi.ID(idSesi), itemByCode.IdSesi)

	items, err := repo.GetSesi(ctx, query.ListSesiFilter{Search: "Integration", Limit: 10})
	require.NoError(t, err)
	assert.Len(t, items, 1)

	updatedCode := "S-INT-02"
	updatedName := "Sesi Updated"
	err = repo.UpdateSesi(ctx, idSesi, updatepatch.UpdateSesiPatch{
		KodeSesi: &updatedCode,
		NamaSesi: &updatedName,
	})
	require.NoError(t, err)

	updated, err := repo.GetSesiById(ctx, idSesi)
	require.NoError(t, err)
	assert.Equal(t, updatedCode, updated.KodeSesi)
	assert.Equal(t, updatedName, updated.NamaSesi)

	err = repo.DeleteSesi(ctx, idSesi)
	require.NoError(t, err)

	_, err = repo.GetSesiById(ctx, idSesi)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
