package integration_test

import (
	"fmt"
	"testing"
	"time"

	aktivitasuserrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAktivitasUserService_CreateAndGet(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := aktivitasuserrepo.NewAktivitasUserRepo(tx, nil)
	svc := aktivitas_user_service.NewAktivitasUserService(repo)

	idPengguna := fixtures.LookupSeedUserID("myadmin")
	description := fmt.Sprintf("integration-aktivitas-%d", time.Now().UnixNano())

	err := svc.CreateAktivitasUserService(ctx, aktivitas_user_service.AktivitasUserCmd{
		IdPengguna:  idPengguna,
		Action:      aktivitas_user.LOGIN,
		Description: "  " + description + "  ",
		IpAddress:   " 127.0.0.1 ",
	})
	require.NoError(t, err)

	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM aktivitas_user
		WHERE id_pengguna = $1
		  AND action = $2
		  AND description = $3
		  AND ip_address = $4
	`, idPengguna, aktivitas_user.LOGIN, description, "127.0.0.1").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	items, err := svc.GetAktivitasUserService(ctx)
	require.NoError(t, err)

	found := false
	for _, item := range items {
		if item.Description == description && item.IpAddress == "127.0.0.1" && item.Action == aktivitas_user.LOGIN {
			found = true
			assert.Equal(t, idPengguna, item.IdPengguna)
			assert.Equal(t, "myadmin", item.Username)
			break
		}
	}
	assert.True(t, found, "expected inserted aktivitas user to be returned by service")
}

func TestAktivitasUserService_CreateAktivitasUserService_InvalidInputDoesNotInsert(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := aktivitasuserrepo.NewAktivitasUserRepo(tx, nil)
	svc := aktivitas_user_service.NewAktivitasUserService(repo)

	idPengguna := fixtures.LookupSeedUserID("myadmin")
	description := fmt.Sprintf("integration-aktivitas-invalid-%d", time.Now().UnixNano())

	before := fixtures.CountRows("aktivitas_user", "description = $1", description)
	require.Equal(t, 0, before)

	err := svc.CreateAktivitasUserService(ctx, aktivitas_user_service.AktivitasUserCmd{
		IdPengguna:  idPengguna,
		Action:      aktivitas_user.Action("INVALID"),
		Description: description,
		IpAddress:   "127.0.0.1",
	})
	require.Error(t, err)

	after := fixtures.CountRows("aktivitas_user", "description = $1", description)
	assert.Equal(t, 0, after)
}
