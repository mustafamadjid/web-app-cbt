package delete_file_repo

import "context"

type DeleteFileRepo interface {
	DeleteFile(ctx context.Context, filePath string) error
}