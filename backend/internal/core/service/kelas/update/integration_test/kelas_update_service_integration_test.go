package integration_test

import (
	"testing"

	kelasrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	kelascreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	kelasupdate "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateKelasService_UpdateNamaKelas(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := kelasrepo.NewKelasRepo(tx, nil)
	createSvc := kelascreate.NewCreateKelasService(repo)
	updateSvc := kelasupdate.NewUpdateKelasService(repo)

	tingkat := 8900
	require.NoError(t, createSvc.CreateTingkatKelas(ctx, kelascreate.CreateTingkatKelasCmd{TingkatKelas: tingkat}))
	var idKelas int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_kelas FROM kelas WHERE tingkat_kelas = $1`, tingkat).Scan(&idKelas))
	require.NoError(t, createSvc.CreateNamaKelas(ctx, kelascreate.CreateNamaKelasCmd{IdTingkatKelas: kelas.ID(idKelas), NamaKelas: "Kelas Lama"}))
	var idNamaKelas int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_nama_kelas FROM nama_kelas WHERE id_kelas = $1 AND nama_kelas = $2`, idKelas, "Kelas Lama").Scan(&idNamaKelas))

	updatedName := "Kelas Baru"
	require.NoError(t, updateSvc.UpdateNamaKelas(ctx, idNamaKelas, updatepatch.NamaKelasPatch{
		IdTingkatKelas: func() *kelas.ID { v := kelas.ID(idKelas); return &v }(),
		NamaKelas:      &updatedName,
	}))

	var stored string
	require.NoError(t, tx.QueryRow(ctx, `SELECT nama_kelas FROM nama_kelas WHERE id_nama_kelas = $1`, idNamaKelas).Scan(&stored))
	assert.Equal(t, updatedName, stored)
}

func TestUpdateKelasService_InvalidInputs(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := kelasrepo.NewKelasRepo(tx, nil)
	updateSvc := kelasupdate.NewUpdateKelasService(repo)

	err := updateSvc.UpdateNamaKelas(ctx, 0, updatepatch.NamaKelasPatch{})
	assert.ErrorIs(t, err, coreerror.ErrMissingId)

	err = updateSvc.UpdateNamaKelas(ctx, 1, updatepatch.NamaKelasPatch{})
	assert.ErrorIs(t, err, coreerror.ErrNoFieldToUpdate)
}

