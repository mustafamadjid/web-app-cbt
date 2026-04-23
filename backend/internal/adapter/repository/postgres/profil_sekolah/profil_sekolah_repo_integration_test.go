package profilsekolahrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfilSekolahRepo_GetAndUpdateProfilSekolah(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewProfilSekolahRepo(tx, nil)

	before, err := repo.GetProfilSekolah(ctx)
	require.NoError(t, err)
	assert.Equal(t, profil_sekolah.IDProfil(1), before.IDProfil)

	logo := "https://example.com/logo.png"
	err = repo.UpdateProfilSekolah(ctx, 1, profil_sekolah.ProfilSekolah{
		EmailSekolah:  "integration-school@example.com",
		NoTelpSekolah: "0211234567",
		KepalaSekolah: "Kepala Integration",
		WakaSekolah:   "Waka Integration",
		NamaSekolah:   "SMA Integration",
		AlamatSekolah: "Jl. Integration 123",
		LogoSekolah:   &logo,
	})
	require.NoError(t, err)

	after, err := repo.GetProfilSekolah(ctx)
	require.NoError(t, err)
	assert.Equal(t, "integration-school@example.com", after.EmailSekolah)
	assert.Equal(t, "0211234567", after.NoTelpSekolah)
	assert.Equal(t, "Kepala Integration", after.KepalaSekolah)
	assert.Equal(t, "Waka Integration", after.WakaSekolah)
	assert.Equal(t, "SMA Integration", after.NamaSekolah)
	assert.Equal(t, "Jl. Integration 123", after.AlamatSekolah)
	require.NotNil(t, after.LogoSekolah)
	assert.Equal(t, logo, *after.LogoSekolah)
}
