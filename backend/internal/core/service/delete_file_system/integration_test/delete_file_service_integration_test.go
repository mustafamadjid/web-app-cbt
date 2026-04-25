package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	delete_file_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/delete_file_system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteFileService_DeleteFile(t *testing.T) {
	t.Run("hapus file berhasil", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		target := filepath.Join(tmpDir, "report.pdf")
		require.NoError(t, os.WriteFile(target, []byte("content"), 0o644))

		svc := delete_file_service.NewDeleteFileService(tmpDir)
		err := svc.DeleteFile(context.Background(), "/uploads/report.pdf")

		require.NoError(t, err)
		_, statErr := os.Stat(target)
		assert.Error(t, statErr)
		assert.True(t, os.IsNotExist(statErr))
	})

	t.Run("path traversal ditolak", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		svc := delete_file_service.NewDeleteFileService(tmpDir)

		err := svc.DeleteFile(context.Background(), "/uploads/../../etc/passwd")

		assert.ErrorIs(t, err, coreerror.ErrPathTraversal)
	})

	t.Run("file tidak ada mengembalikan error os", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		svc := delete_file_service.NewDeleteFileService(tmpDir)

		err := svc.DeleteFile(context.Background(), "/uploads/missing.txt")

		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})
}

