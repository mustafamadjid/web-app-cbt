package pengumumanrepo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPengumumanRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewPengumumanRepo(tx, nil)

	createdUser := fixtures.CreateUser(user.ADMIN)
	now := time.Now().UTC()
	active := fixtures.CreatePengumuman(createdUser.ID, now.AddDate(0, 0, -1), now.AddDate(0, 0, 1))
	nonActive := fixtures.CreatePengumuman(createdUser.ID, now.AddDate(0, 0, -10), now.AddDate(0, 0, -1))
	incoming := fixtures.CreatePengumuman(createdUser.ID, now.AddDate(0, 0, 1), now.AddDate(0, 0, 2))

	activeItems, err := repo.GetPengumumanActive(ctx)
	require.NoError(t, err)
	assert.Contains(t, extractPengumumanIDs(activeItems), pengumuman.ID(active.ID))

	nonActiveItems, err := repo.GetPengumumanNonActive(ctx)
	require.NoError(t, err)
	assert.Contains(t, extractPengumumanIDs(nonActiveItems), pengumuman.ID(nonActive.ID))

	incomingItems, err := repo.GetPengumumanIncoming(ctx)
	require.NoError(t, err)
	assert.Contains(t, extractPengumumanIDs(incomingItems), pengumuman.ID(incoming.ID))

	createdTitle := "Repo Announcement"
	err = repo.CreatePengumuman(ctx, pengumuman.Pengumuman{
		IdPengguna:               pengumuman.ID(createdUser.ID),
		JudulPengumuman:          createdTitle,
		IsiPengumuman:            "isi repo announcement",
		TanggalRilisPengumuman:   now.Format("2006-01-02"),
		TanggalSelesaiPengumuman: now.AddDate(0, 0, 3).Format("2006-01-02"),
		DokumenPengumuman:        "",
	})
	require.NoError(t, err)

	var createdID int
	err = tx.QueryRow(ctx, `SELECT id_pengumuman FROM pengumuman WHERE judul_pengumuman = $1`, createdTitle).Scan(&createdID)
	require.NoError(t, err)

	item, err := repo.GetPengumumanById(ctx, pengumuman.ID(createdID))
	require.NoError(t, err)
	assert.Equal(t, createdTitle, item.JudulPengumuman)

	updatedTitle := "Repo Announcement Updated"
	err = repo.UpdatePengumuman(ctx, pengumuman.ID(createdID), updatepatch.PengumumanUpdatePatch{
		JudulPengumuman: &updatedTitle,
	})
	require.NoError(t, err)

	updated, err := repo.GetPengumumanById(ctx, pengumuman.ID(createdID))
	require.NoError(t, err)
	assert.Equal(t, updatedTitle, updated.JudulPengumuman)

	err = repo.DeletePengumuman(ctx, pengumuman.ID(createdID))
	require.NoError(t, err)

	_, err = repo.GetPengumumanById(ctx, pengumuman.ID(createdID))
	assert.ErrorIs(t, err, coreerror.ErrNotFound)

	_, err = repo.GetPengumumanById(ctx, pengumuman.ID(999999))
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

func extractPengumumanIDs(items []pengumuman.Pengumuman) []pengumuman.ID {
	ids := make([]pengumuman.ID, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.IdPengumuman)
	}
	return ids
}
