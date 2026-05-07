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

func TestDeleteFileService_BranchCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		filePath     string
		setupFile    string
		wantErr      error
		wantNotExist bool
	}{
		{
			name:     "Branch 1 -> path traversal ditolak",
			filePath: "/uploads/../../etc/passwd",
			wantErr:  coreerror.ErrPathTraversal,
		},
		{
			name:         "Branch 2 -> file tidak ada mengembalikan error",
			filePath:     "/uploads/nonexistent.txt",
			wantNotExist: true,
		},
		{
			name:      "Branch 3 -> file valid berhasil dihapus",
			filePath:  "/uploads/testfile.txt",
			setupFile: "testfile.txt",
		},
		{
			name:      "Branch 4 -> path dengan prefix uploads dibersihkan dengan benar",
			filePath:  "/uploads/nested/final.txt",
			setupFile: filepath.Join("nested", "final.txt"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			if tc.setupFile != "" {
				testFile := filepath.Join(tmpDir, tc.setupFile)
				err := os.MkdirAll(filepath.Dir(testFile), 0o755)
				assert.NoError(t, err)
				err = os.WriteFile(testFile, []byte("test content"), 0o644)
				assert.NoError(t, err)
			}

			svc := delete_file_service.NewDeleteFileService(tmpDir)

			err := svc.DeleteFile(context.Background(), tc.filePath)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			if tc.wantNotExist {
				assert.Error(t, err)
				assert.True(t, os.IsNotExist(err))
				return
			}

			assert.NoError(t, err)

			deletedPath := filepath.Join(tmpDir, filepath.Clean(filepath.FromSlash(tc.setupFile)))
			_, statErr := os.Stat(deletedPath)
			assert.True(t, os.IsNotExist(statErr))
		})
	}
}
