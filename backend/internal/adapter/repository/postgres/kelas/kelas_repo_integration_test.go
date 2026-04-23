package kelasrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKelasRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewKelasRepo(tx, nil)

	err := repo.CreateTingkatKelas(ctx, 88)
	require.NoError(t, err)

	var idKelas int
	err = tx.QueryRow(ctx, `SELECT id_kelas FROM kelas WHERE tingkat_kelas = 88`).Scan(&idKelas)
	require.NoError(t, err)

	exists, err := repo.ExistTingkatKelas(ctx, 88)
	require.NoError(t, err)
	assert.True(t, exists)

	err = repo.CreateNamaKelas(ctx, kelas.NamaKelas{
		IdTingkatKelas: kelas.ID(idKelas),
		NamaKelas:      "Kelas IT 88",
	})
	require.NoError(t, err)

	var idNamaKelas int
	err = tx.QueryRow(ctx, `SELECT id_nama_kelas FROM nama_kelas WHERE id_kelas = $1 AND nama_kelas = $2`, idKelas, "Kelas IT 88").Scan(&idNamaKelas)
	require.NoError(t, err)

	nameExists, err := repo.ExistNamaKelas(ctx, "Kelas IT 88")
	require.NoError(t, err)
	assert.True(t, nameExists)

	items, err := repo.GetKelas(ctx, query.ListKelasFilter{Search: "IT 88", Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.NotEmpty(t, items[0].ItemsTingkatKelas)
	assert.NotEmpty(t, items[0].ItemsNamaKelas)

	item, err := repo.GetKelasById(ctx, idKelas, idNamaKelas)
	require.NoError(t, err)
	assert.Equal(t, 88, item.ItemsTingkatKelas.TingkatKelas)
	assert.Equal(t, "Kelas IT 88", item.ItemsNamaKelas.NamaKelas)

	updatedName := "Kelas IT 88 Updated"
	newTingkat := kelas.ID(idKelas)
	err = repo.UpdateNamaKelas(ctx, idNamaKelas, updatepatch.NamaKelasPatch{
		IdTingkatKelas: &newTingkat,
		NamaKelas:      &updatedName,
	})
	require.NoError(t, err)

	updated, err := repo.GetKelasById(ctx, idKelas, idNamaKelas)
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.ItemsNamaKelas.NamaKelas)

	err = repo.DeleteNamaKelas(ctx, idNamaKelas)
	require.NoError(t, err)

	_, err = repo.GetKelasById(ctx, idKelas, idNamaKelas)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

func TestKelasRepo_DeleteNamaKelas_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewKelasRepo(tx, nil)

	err := repo.DeleteNamaKelas(ctx, 999999)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
