package bank_soal_service_test

import (
	"context"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/update"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/update/fake_test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBankSoalService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &fake_test.FakeBankSoalRepo{}
	svc := bank_soal_service.NewUpdateBankSoalService(repo)

	idMapel := bank_soal.ID(2)
	idKelas := bank_soal.ID(3)
	idPengguna := bank_soal.ID(4)

	patch := updatepatch.UpdateBankSoalPatch{
		IdMapel:      &idMapel,
		IdKelas:      &idKelas,
		IdPengguna:   &idPengguna,
		NamaBankSoal: ptrString("  Nama Bank Soal  "),
		Deskripsi:    ptrString("  Deskripsi  "),
		Materi:       ptrString("  Materi  "),
	}

	err := svc.UpdateBankSoalService(ctx, 10, patch)
	require.NoError(t, err)
	require.True(t, repo.UpdateCalled)
	assert.Equal(t, bank_soal.ID(10), repo.GotID)

	require.NotNil(t, repo.GotPatch.NamaBankSoal)
	require.NotNil(t, repo.GotPatch.Deskripsi)
	require.NotNil(t, repo.GotPatch.Materi)
	assert.Equal(t, "Nama Bank Soal", *repo.GotPatch.NamaBankSoal)
	assert.Equal(t, "Deskripsi", *repo.GotPatch.Deskripsi)
	assert.Equal(t, "Materi", *repo.GotPatch.Materi)

	require.NotNil(t, repo.GotPatch.IdMapel)
	require.NotNil(t, repo.GotPatch.IdKelas)
	require.NotNil(t, repo.GotPatch.IdPengguna)
	assert.Equal(t, idMapel, *repo.GotPatch.IdMapel)
	assert.Equal(t, idKelas, *repo.GotPatch.IdKelas)
	assert.Equal(t, idPengguna, *repo.GotPatch.IdPengguna)
}
