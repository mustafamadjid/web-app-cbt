package delete_file_service_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	delete_file_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/delete_file_system"
	"github.com/stretchr/testify/assert"
)

func TestDeleteFileService_BasisPath(t *testing.T) {
	t.Parallel()

	t.Run("Path 1 -> path traversal detected", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		svc := delete_file_service.NewDeleteFileService(tmpDir)

		ctx := context.Background()
		err := svc.DeleteFile(ctx, "/uploads/../../etc/passwd")

		assert.ErrorIs(t, err, coreerror.ErrPathTraversal)
	})

	t.Run("Path 2 -> os.Remove error (file tidak ada)", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		svc := delete_file_service.NewDeleteFileService(tmpDir)

		ctx := context.Background()
		err := svc.DeleteFile(ctx, "/uploads/nonexistent.txt")

		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Path 3 -> happy path berhasil delete file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		// Buat file sementara untuk di-delete
		testFile := filepath.Join(tmpDir, "testfile.txt")
		err := os.WriteFile(testFile, []byte("test content"), 0o644)
		assert.NoError(t, err)

		svc := delete_file_service.NewDeleteFileService(tmpDir)

		ctx := context.Background()
		err = svc.DeleteFile(ctx, "/uploads/testfile.txt")

		assert.NoError(t, err)

		// Pastikan file sudah terhapus
		_, statErr := os.Stat(testFile)
		assert.True(t, os.IsNotExist(statErr))
	})
}
