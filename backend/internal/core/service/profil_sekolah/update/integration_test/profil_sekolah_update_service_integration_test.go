package integration_test

import (
	"testing"

	profilsekolahrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_sekolah"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	profilsekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfilSekolahService_Update(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := profilsekolahrepo.NewProfilSekolahRepo(tx, nil)
	svc := profilsekolah_service.NewUpdateProfilSekolahService(repo)

	email := "integration-school@example.com"
	noTelp := "0211234567"
	kepala := "Kepala Integration"
	waka := "Waka Integration"
	nama := "SMA Integration"
	alamat := "Jl. Integration 123"
	logo := "https://example.com/logo.png"
	require.NoError(t, svc.UpdateProfilSekolah(ctx, profilsekolah_service.UpdateProfilSekolahCmd{
		IDProfil:      1,
		EmailSekolah:  &email,
		NoTelpSekolah: &noTelp,
		KepalaSekolah: &kepala,
		WakaSekolah:   &waka,
		NamaSekolah:   &nama,
		AlamatSekolah: &alamat,
		LogoSekolah:   &logo,
	}))

	item, err := repo.GetProfilSekolah(ctx)
	require.NoError(t, err)
	assert.Equal(t, email, item.EmailSekolah)
	assert.Equal(t, noTelp, item.NoTelpSekolah)
	assert.Equal(t, kepala, item.KepalaSekolah)
	assert.Equal(t, waka, item.WakaSekolah)
	assert.Equal(t, nama, item.NamaSekolah)
	assert.Equal(t, alamat, item.AlamatSekolah)
	require.NotNil(t, item.LogoSekolah)
	assert.Equal(t, logo, *item.LogoSekolah)
}

