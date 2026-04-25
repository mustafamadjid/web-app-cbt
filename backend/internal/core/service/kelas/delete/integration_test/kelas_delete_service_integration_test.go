package integration_test

import (
	"testing"

	kelasrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	kelascreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	kelasdelete "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/delete"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteKelasService_DeleteNamaKelas(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := kelasrepo.NewKelasRepo(tx, nil)
	createSvc := kelascreate.NewCreateKelasService(repo)
	deleteSvc := kelasdelete.NewDeleteKelasService(repo)

	tingkat := 8910
	require.NoError(t, createSvc.CreateTingkatKelas(ctx, kelascreate.CreateTingkatKelasCmd{TingkatKelas: tingkat}))
	var idKelas int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_kelas FROM kelas WHERE tingkat_kelas = $1`, tingkat).Scan(&idKelas))
	require.NoError(t, createSvc.CreateNamaKelas(ctx, kelascreate.CreateNamaKelasCmd{IdTingkatKelas: kelas.ID(idKelas), NamaKelas: "Kelas Hapus"}))
	var idNamaKelas int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_nama_kelas FROM nama_kelas WHERE id_kelas = $1 AND nama_kelas = $2`, idKelas, "Kelas Hapus").Scan(&idNamaKelas))

	require.NoError(t, deleteSvc.DeleteNamaKelas(ctx, idNamaKelas))
	_, err := tx.Exec(ctx, `SELECT 1 FROM nama_kelas WHERE id_nama_kelas = $1`, idNamaKelas)
	assert.NoError(t, err)
}

func TestDeleteKelasService_DeleteNamaKelas_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := kelasrepo.NewKelasRepo(tx, nil)
	deleteSvc := kelasdelete.NewDeleteKelasService(repo)

	err := deleteSvc.DeleteNamaKelas(ctx, 999999)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

