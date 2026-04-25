package integration_test

import (
	"time"
	"testing"

	sessionrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/session"
	sesirepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/sesi"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	sesisdomain "github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	sesicreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
	sesiget "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/get"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSesiService_GetAndActiveSessions(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := sesirepo.NewSesirepo(tx, nil)
	sessionRepo := sessionrepo.NewSessionRepo(tx, nil)
	createSvc := sesicreate.NewCreateSesiService(repo)
	getSvc := sesiget.NewGetSesiService(repo, sessionRepo)

	require.NoError(t, createSvc.CreateSesiService(ctx, sesisdomain.Sesi{KodeSesi: "S-GIT-01", NamaSesi: "Sesi Get Integration"}))
	items, err := getSvc.GetSesiService(ctx, query.ListSesiFilter{Search: "Get Integration", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, items)

	var idSesi int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_sesi FROM sesi_ujian WHERE kode_sesi = 'S-GIT-01'`).Scan(&idSesi))
	item, err := getSvc.GetSesiByIdService(ctx, idSesi)
	require.NoError(t, err)
	assert.Equal(t, "Sesi Get Integration", item.NamaSesi)

	itemByCode, err := getSvc.GetSesiByKodeService(ctx, "S-GIT-01")
	require.NoError(t, err)
	assert.Equal(t, sesisdomain.ID(idSesi), itemByCode.IdSesi)

	u := fixtures.CreateUser(user.ADMIN)
	sessID, err := sessionRepo.CreateSession(ctx, u.ID, u.Role, time.Now().Add(time.Hour))
	require.NoError(t, err)

	active, err := getSvc.GetAllActiveSessionService(ctx)
	require.NoError(t, err)
	found := false
	for _, item := range active {
		if item.Session.SessionID == sessID {
			found = true
		}
	}
	assert.True(t, found)
}

func TestGetSesiService_InvalidID(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := sesirepo.NewSesirepo(tx, nil)
	sessionRepo := sessionrepo.NewSessionRepo(tx, nil)
	getSvc := sesiget.NewGetSesiService(repo, sessionRepo)

	_, err := getSvc.GetSesiByIdService(ctx, 0)
	assert.ErrorIs(t, err, coreerror.ErrMissingId)
}
