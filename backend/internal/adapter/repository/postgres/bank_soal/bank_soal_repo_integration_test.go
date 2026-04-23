package banksoalrepo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	banksoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBankSoalRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewBankSoalRepo(tx, nil)

	kelas := fixtures.CreateKelas(71)
	guru := fixtures.CreateUser("GURU")
	fixtures.CreateGuruProfile(guru.ID)
	mapel := fixtures.CreateMapel(kelas.ID)

	err := repo.CreateBankSoal(ctx, banksoal.BankSoal{
		IdMapel:      banksoal.ID(mapel.ID),
		IdKelas:      banksoal.ID(kelas.ID),
		IdPengguna:   banksoal.ID(guru.ID),
		NamaBankSoal: "Bank Repo",
		Deskripsi:    "Desc Repo",
		Materi:       "Materi Repo",
	})
	require.NoError(t, err)

	var bankID int
	err = tx.QueryRow(ctx, `SELECT id_bank_soal FROM bank_soal WHERE nama_bank_soal = 'Bank Repo'`).Scan(&bankID)
	require.NoError(t, err)

	item, err := repo.GetBankSoalById(ctx, banksoal.ID(bankID))
	require.NoError(t, err)
	assert.Equal(t, "Bank Repo", item.NamaBankSoal)
	assert.False(t, item.SoalUploaded)

	items, err := repo.GetBankSoal(ctx, query.BankSoalFilter{Search: "Bank Repo", Limit: 10})
	require.NoError(t, err)
	assert.Len(t, items, 1)

	byGuru, err := repo.GetBankSoalByGuru(ctx, banksoal.ID(guru.ID))
	require.NoError(t, err)
	assert.Len(t, byGuru, 1)

	version := fixtures.CreateBankSoalVersion(int64(bankID), guru.ID, "published", 1)
	fixtures.SetActiveBankSoalVersion(int64(bankID), version.ID)

	uploaded, err := repo.GetBankSoalUploaded(ctx, query.BankSoalFilter{Search: "Bank Repo", Limit: 10})
	require.NoError(t, err)
	assert.Len(t, uploaded, 1)
	assert.True(t, uploaded[0].SoalUploaded)

	updatedName := "Bank Repo Updated"
	updatedMateri := "Materi Updated"
	err = repo.UpdateBankSoal(ctx, banksoal.ID(bankID), updatepatch.UpdateBankSoalPatch{
		NamaBankSoal: &updatedName,
		Materi:       &updatedMateri,
	})
	require.NoError(t, err)

	updated, err := repo.GetBankSoalById(ctx, banksoal.ID(bankID))
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.NamaBankSoal)
	assert.Equal(t, updatedMateri, updated.Materi)

	err = repo.DeleteBankSoal(ctx, banksoal.ID(bankID))
	require.NoError(t, err)

	_, err = repo.GetBankSoalById(ctx, banksoal.ID(bankID))
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

func TestBankSoalRepo_GetIdBankSoalByAttemptId(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewBankSoalRepo(tx, nil)

	exam := fixtures.CreateExamFixture()
	start := time.Date(2099, time.January, 1, 8, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_IN_PROGRESS, &start, nil, testutil.Ptr(start.Add(time.Hour)))

	idBankSoal, err := repo.GetIdBankSoalByAttemptId(ctx, ujian.ID(attempt.ID))
	require.NoError(t, err)
	assert.Equal(t, ujian.ID(exam.BankSoal.ID), idBankSoal)
}
