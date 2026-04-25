package integration_test

import (
	"testing"

	profilsekolahrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_sekolah"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	profilsekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfilSekolahService_Get(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := profilsekolahrepo.NewProfilSekolahRepo(tx, nil)
	svc := profilsekolah_service.NewGetProfilSekolahService(repo)

	item, err := svc.GetProfilSekolah(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, int(item.IDProfil))
}
