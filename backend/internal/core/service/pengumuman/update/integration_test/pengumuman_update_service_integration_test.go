package integration_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	deletefilesvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/delete_file_system"
	pengumumanrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumumancreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/create"
	pengumumanupdate "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPengumumanService_UpdateAndDeleteDocument(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := pengumumanrepo.NewPengumumanRepo(tx, nil)
	tmpDir := t.TempDir()
	deleteFileSvc := deletefilesvc.NewDeleteFileService(tmpDir)
	createSvc := pengumumancreate.NewCreatePengumumanRepo(repo)
	updateSvc := pengumumanupdate.NewUpdatePengumumanService(repo, deleteFileSvc)

	creator := fixtures.CreateUser(user.ADMIN)
	oldDoc := "old-announcement.pdf"
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, oldDoc), []byte("old"), 0o644))

	now := time.Now().UTC()
	require.NoError(t, createSvc.CreatePengumuman(ctx, pengumuman.Pengumuman{
		IdPengguna:               pengumuman.ID(creator.ID),
		JudulPengumuman:          "Pengumuman Lama",
		IsiPengumuman:            "Isi lama",
		TanggalRilisPengumuman:   now.Add(-48 * time.Hour).Format("2006-01-02"),
		TanggalSelesaiPengumuman: now.Add(48 * time.Hour).Format("2006-01-02"),
		DokumenPengumuman:        "/uploads/" + oldDoc,
	}))

	var idPengumuman int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_pengumuman FROM pengumuman WHERE judul_pengumuman = 'Pengumuman Lama'`).Scan(&idPengumuman))

	newDoc := "new-announcement.pdf"
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, newDoc), []byte("new"), 0o644))

	newTitle := " Pengumuman Baru "
	require.NoError(t, updateSvc.UpdatePengumumanService(ctx, pengumuman.ID(idPengumuman), updatepatch.PengumumanUpdatePatch{
		JudulPengumuman:   &newTitle,
		DokumenPengumuman: func() *string { v := "/uploads/" + newDoc; return &v }(),
	}))

	item, err := repo.GetPengumumanById(ctx, pengumuman.ID(idPengumuman))
	require.NoError(t, err)
	assert.Equal(t, "Pengumuman Baru", item.JudulPengumuman)

	_, err = os.Stat(filepath.Join(tmpDir, oldDoc))
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

