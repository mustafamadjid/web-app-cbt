package integration_test

import (
	"testing"
	"time"

	pengumumanrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	pengumumancreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/create"
	pengumumanget "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPengumumanService_CreateAndGet(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := pengumumanrepo.NewPengumumanRepo(tx, nil)
	createSvc := pengumumancreate.NewCreatePengumumanRepo(repo)
	getSvc := pengumumanget.NewGetPengumumanService(repo)

	creator := fixtures.CreateUser(user.ADMIN)
	now := time.Now().UTC()
	title := "  Pengumuman Integration  "
	require.NoError(t, createSvc.CreatePengumuman(ctx, pengumuman.Pengumuman{
		IdPengguna:               pengumuman.ID(creator.ID),
		JudulPengumuman:          title,
		IsiPengumuman:            "  Isi integration  ",
		TanggalRilisPengumuman:   now.Add(-time.Hour).Format("2006-01-02"),
		TanggalSelesaiPengumuman: now.Add(time.Hour).Format("2006-01-02"),
		DokumenPengumuman:        "",
	}))

	active, err := getSvc.GetPengumumanActiveService(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, active)

	found, err := getSvc.GetPengumumanByIdService(ctx, active[0].IdPengumuman)
	require.NoError(t, err)
	assert.Equal(t, "Pengumuman Integration", found.JudulPengumuman)
}

func TestPengumumanService_CreateInvalidDate(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := pengumumanrepo.NewPengumumanRepo(tx, nil)
	createSvc := pengumumancreate.NewCreatePengumumanRepo(repo)

	err := createSvc.CreatePengumuman(ctx, pengumuman.Pengumuman{IdPengguna: 1, JudulPengumuman: "A", IsiPengumuman: "B", TanggalRilisPengumuman: "2026-13-01", TanggalSelesaiPengumuman: "2026-13-02"})
	assert.ErrorIs(t, err, coreerror.ErrInvalidDateFormat)
}
