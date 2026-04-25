package integration_test

import (
	"testing"

	sesirepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/sesi"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	sesicreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
	sesidelete "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/delete"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteSesiService_Delete(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := sesirepo.NewSesirepo(tx, nil)
	createSvc := sesicreate.NewCreateSesiService(repo)
	deleteSvc := sesidelete.NewDeleteSesiService(repo)

	require.NoError(t, createSvc.CreateSesiService(ctx, sesi.Sesi{KodeSesi: "S-DEL-01", NamaSesi: "Sesi Delete"}))
	var idSesi int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_sesi FROM sesi_ujian WHERE kode_sesi = 'S-DEL-01'`).Scan(&idSesi))

	require.NoError(t, deleteSvc.DeleteSesiService(ctx, idSesi))
	_, err := repo.GetSesiById(ctx, idSesi)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

