package integration_test

import (
	"testing"

	kelasrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	kelascreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	kelasget "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetKelasService_GetFullAndByID(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := kelasrepo.NewKelasRepo(tx, nil)
	createSvc := kelascreate.NewCreateKelasService(repo)
	getSvc := kelasget.NewGetKelasService(repo)

	tingkat := 8800
	require.NoError(t, createSvc.CreateTingkatKelas(ctx, kelascreate.CreateTingkatKelasCmd{TingkatKelas: tingkat}))
	var idKelas int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_kelas FROM kelas WHERE tingkat_kelas = $1`, tingkat).Scan(&idKelas))

	nama := "Kelas Integration"
	require.NoError(t, createSvc.CreateNamaKelas(ctx, kelascreate.CreateNamaKelasCmd{IdTingkatKelas: kelas.ID(idKelas), NamaKelas: nama}))
	var idNamaKelas int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_nama_kelas FROM nama_kelas WHERE id_kelas = $1 AND nama_kelas = $2`, idKelas, nama).Scan(&idNamaKelas))

	items, err := getSvc.GetFullKelas(ctx, query.ListKelasFilter{Search: "Integration", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, items)

	found := false
	for _, item := range items {
		if len(item.ItemsNamaKelas) > 0 && item.ItemsNamaKelas[0].NamaKelas == nama {
			found = true
			require.NotEmpty(t, item.ItemsTingkatKelas)
			assert.Equal(t, tingkat, item.ItemsTingkatKelas[0].TingkatKelas)
		}
	}
	assert.True(t, found)

	item, err := getSvc.GetKelasById(ctx, idKelas, idNamaKelas)
	require.NoError(t, err)
	assert.Equal(t, tingkat, item.ItemsTingkatKelas.TingkatKelas)
	assert.Equal(t, nama, item.ItemsNamaKelas.NamaKelas)

	_ = fixtures
}

func TestGetKelasService_GetKelasByID_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := kelasrepo.NewKelasRepo(tx, nil)
	getSvc := kelasget.NewGetKelasService(repo)

	_, err := getSvc.GetKelasById(ctx, 999999, 999999)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
