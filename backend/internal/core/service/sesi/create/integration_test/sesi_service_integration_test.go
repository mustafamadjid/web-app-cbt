package integration_test

import (
	"testing"

	sesirepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/sesi"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	sesisdomain "github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	sesisvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSesiService_Create(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := sesirepo.NewSesirepo(tx, nil)
	svc := sesisvc.NewCreateSesiService(repo)

	err := svc.CreateSesiService(ctx, sesisdomain.Sesi{KodeSesi: " s-int-01 ", NamaSesi: " Sesi Integration "})
	require.NoError(t, err)

	items, err := repo.GetSesi(ctx, query.ListSesiFilter{Search: "Integration", Limit: 10})
	require.NoError(t, err)
	assert.NotEmpty(t, items)
}

func TestCreateSesiService_DuplicateAndInvalid(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := sesirepo.NewSesirepo(tx, nil)
	svc := sesisvc.NewCreateSesiService(repo)

	require.NoError(t, svc.CreateSesiService(ctx, sesisdomain.Sesi{KodeSesi: "S-INT-02", NamaSesi: "Sesi 2"}))
	err := svc.CreateSesiService(ctx, sesisdomain.Sesi{KodeSesi: "S-INT-02", NamaSesi: "Sesi 2 Lain"})
	assert.ErrorIs(t, err, coreerror.ErrSesiUjianExist)

	err = svc.CreateSesiService(ctx, sesisdomain.Sesi{})
	assert.ErrorIs(t, err, coreerror.ErrMissingField)
}
