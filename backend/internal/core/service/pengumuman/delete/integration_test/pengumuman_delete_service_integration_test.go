package integration_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	deletefilesvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/delete_file_system"
	pengumumanrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	pengumumancreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/create"
	pengumumandelete "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/delete"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPengumumanService_Delete(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := pengumumanrepo.NewPengumumanRepo(tx, nil)
	tmpDir := t.TempDir()
	deleteFileSvc := deletefilesvc.NewDeleteFileService(tmpDir)
	createSvc := pengumumancreate.NewCreatePengumumanRepo(repo)
	deleteSvc := pengumumandelete.NewDeletePengumumanService(repo, deleteFileSvc)

	creator := fixtures.CreateUser(user.ADMIN)
	doc := "delete-announcement.pdf"
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, doc), []byte("doc"), 0o644))
	now := time.Now().UTC()
	require.NoError(t, createSvc.CreatePengumuman(ctx, pengumuman.Pengumuman{
		IdPengguna:               pengumuman.ID(creator.ID),
		JudulPengumuman:          "Pengumuman Hapus",
		IsiPengumuman:            "Isi hapus",
		TanggalRilisPengumuman:   now.Add(-48 * time.Hour).Format("2006-01-02"),
		TanggalSelesaiPengumuman: now.Add(48 * time.Hour).Format("2006-01-02"),
		DokumenPengumuman:        "/uploads/" + doc,
	}))

	var idPengumuman int
	require.NoError(t, tx.QueryRow(ctx, `SELECT id_pengumuman FROM pengumuman WHERE judul_pengumuman = 'Pengumuman Hapus'`).Scan(&idPengumuman))

	require.NoError(t, deleteSvc.DeletePengumumanService(ctx, pengumuman.ID(idPengumuman)))

	_, err := repo.GetPengumumanById(ctx, pengumuman.ID(idPengumuman))
	assert.ErrorIs(t, err, coreerror.ErrNotFound)

	_, err = os.Stat(filepath.Join(tmpDir, doc))
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

