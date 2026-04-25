package integration_test

import (
	"testing"

	banksoalrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/bank_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	bank_soal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	banksoalcreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/create"
	banksoaldelete "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/delete"
	banksoalget "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/get"
	banksoalupdate "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/update"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBankSoalService_CreateGetUpdateDelete(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	kelas := fixtures.CreateKelas(8700)
	mapel := fixtures.CreateMapel(kelas.ID)
	guru := fixtures.CreateUser(user.GURU)

	repo := banksoalrepo.NewBankSoalRepo(tx, nil)
	createSvc := banksoalcreate.NewCreateBankSoalService(repo)
	getSvc := banksoalget.NewGetBankSoalService(repo)
	updateSvc := banksoalupdate.NewUpdateBankSoalService(repo)
	deleteSvc := banksoaldelete.NewDeleteBankSoalService(repo)

	require.NoError(t, createSvc.CreateBankSoalService(ctx, bank_soal.BankSoal{
		IdMapel:      bank_soal.ID(mapel.ID),
		IdKelas:      bank_soal.ID(kelas.ID),
		IdPengguna:   bank_soal.ID(guru.ID),
		NamaBankSoal: "  Bank Integration  ",
		Deskripsi:    "  Deskripsi Bank  ",
		Materi:       "  Materi Bank  ",
	}))

	items, err := getSvc.GetBankSoalService(ctx, query.BankSoalFilter{Search: "Integration", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, items)

	var idBankSoal int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_bank_soal FROM bank_soal WHERE nama_bank_soal = 'Bank Integration'`).Scan(&idBankSoal))

	item, err := getSvc.GetBankSoalByIdService(ctx, bank_soal.ID(idBankSoal))
	require.NoError(t, err)
	assert.Equal(t, "Bank Integration", item.NamaBankSoal)

	byGuru, err := getSvc.GetBankSoalByGuruService(ctx, bank_soal.ID(guru.ID))
	require.NoError(t, err)
	require.NotEmpty(t, byGuru)

	newName := "Bank Baru"
	newDesc := "Deskripsi Baru"
	require.NoError(t, updateSvc.UpdateBankSoalService(ctx, bank_soal.ID(idBankSoal), updatepatch.UpdateBankSoalPatch{
		NamaBankSoal: &newName,
		Deskripsi:    &newDesc,
	}))

	updated, err := getSvc.GetBankSoalByIdService(ctx, bank_soal.ID(idBankSoal))
	require.NoError(t, err)
	assert.Equal(t, newName, updated.NamaBankSoal)
	assert.Equal(t, newDesc, updated.Deskripsi)

	require.NoError(t, deleteSvc.DeleteBankSoalService(ctx, bank_soal.ID(idBankSoal)))
	_, err = getSvc.GetBankSoalByIdService(ctx, bank_soal.ID(idBankSoal))
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
