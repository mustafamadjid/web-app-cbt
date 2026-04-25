package integration_test

import (
	"testing"

	sesirepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/sesi"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesicreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
	sesiupdate "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/update"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateSesiService_Update(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := sesirepo.NewSesirepo(tx, nil)
	createSvc := sesicreate.NewCreateSesiService(repo)
	updateSvc := sesiupdate.NewUpdateSesiService(repo)

	require.NoError(t, createSvc.CreateSesiService(ctx, sesi.Sesi{KodeSesi: "S-UP-01", NamaSesi: "Sesi Update"}))
	var idSesi int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_sesi FROM sesi_ujian WHERE kode_sesi = 'S-UP-01'`).Scan(&idSesi))

	updatedCode := "s-up-02"
	updatedName := "  Sesi Baru  "
	require.NoError(t, updateSvc.UpdateSesiService(ctx, idSesi, updatepatch.UpdateSesiPatch{
		KodeSesi: &updatedCode,
		NamaSesi: &updatedName,
	}))

	item, err := repo.GetSesiById(ctx, idSesi)
	require.NoError(t, err)
	assert.Equal(t, "S-UP-02", item.KodeSesi)
	assert.Equal(t, "Sesi Baru", item.NamaSesi)
}

