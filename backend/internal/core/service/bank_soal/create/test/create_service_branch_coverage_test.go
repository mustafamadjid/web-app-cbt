package bank_soal_service_test

import (
	"context"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBankSoalService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &FakeBankSoalRepo{}
	svc := bank_soal_service.NewCreateBankSoalService(repo)

	input := bank_soal.BankSoal{
		IdMapel:      10,
		IdKelas:      11,
		IdPengguna:   12,
		NamaBankSoal: "  Nama Bank Soal  ",
		Deskripsi:    "  Deskripsi  ",
		Materi:       "  Materi  ",
	}

	err := svc.CreateBankSoalService(ctx, input)
	require.NoError(t, err)
	require.True(t, repo.CreateCalled)

	assert.Equal(t, "Nama Bank Soal", repo.GotCreate.NamaBankSoal)
	assert.Equal(t, "Deskripsi", repo.GotCreate.Deskripsi)
	assert.Equal(t, "Materi", repo.GotCreate.Materi)
	assert.Equal(t, bank_soal.ID(10), repo.GotCreate.IdMapel)
	assert.Equal(t, bank_soal.ID(11), repo.GotCreate.IdKelas)
	assert.Equal(t, bank_soal.ID(12), repo.GotCreate.IdPengguna)
}
