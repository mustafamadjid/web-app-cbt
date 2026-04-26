package importsoaljobrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportSoalJobRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewImportSoalJobRepo(tx, nil)

	kelas := fixtures.CreateKelas(74)
	guru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guru.ID)
	mapel := fixtures.CreateMapel(kelas.ID)
	bank := fixtures.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)

	jobID, err := repo.CreateJob(ctx, importsoal.ImportSoalJob{
		IDBankSoal: bank.ID,
		IDPengguna: int64(guru.ID),
		Status:     importsoal.StatusPending,
		FilePath:   "/tmp/import.docx",
	})
	require.NoError(t, err)

	got, err := repo.GetJobByID(ctx, jobID)
	require.NoError(t, err)
	assert.Equal(t, bank.ID, got.IDBankSoal)
	assert.Equal(t, importsoal.StatusPending, got.Status)

	pending, err := repo.GetPendingJobs(ctx, 10)
	require.NoError(t, err)
	assert.Len(t, pending, 1)

	err = repo.UpdateJobStatus(ctx, jobID, importsoal.StatusCompleted, "", "", 7)
	require.NoError(t, err)

	updated, err := repo.GetJobByID(ctx, jobID)
	require.NoError(t, err)
	assert.Equal(t, importsoal.StatusCompleted, updated.Status)
	assert.Equal(t, 7, updated.TotalSoal)

	items, err := repo.GetJobsByBankSoal(ctx, bank.ID)
	require.NoError(t, err)
	assert.Len(t, items, 1)
}

func TestImportSoalJobRepo_GetJobByID_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewImportSoalJobRepo(tx, nil)

	_, err := repo.GetJobByID(ctx, 999999)
	assert.ErrorIs(t, err, coreerror.ErrImportJobNotFound)
}
