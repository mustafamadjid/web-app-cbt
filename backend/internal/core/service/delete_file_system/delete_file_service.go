package delete_file_service

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type DeleteFileService struct {
	uploadDir string
}

func NewDeleteFileService( uploadDir string) *DeleteFileService {
	return &DeleteFileService{uploadDir: uploadDir}
}

func(r *DeleteFileService)DeleteFile(ctx context.Context, filePath string) error {
	logger := corelog.FromContext(ctx)

	rel := strings.TrimPrefix(filePath,"/uploads/")
	rel = filepath.Clean(filepath.FromSlash(rel))

	target := filepath.Join(r.uploadDir,rel)



	relative,err := filepath.Rel(r.uploadDir,target)
	if err != nil {
		logger.Error(ctx,"failed delete file","layer","core.service","op","delete_file","err",err)
		return err
	}
	if relative == ".." || strings.HasPrefix(relative,".." + string(filepath.Separator)){
		logger.Error(ctx,"failed delete file","layer","core.service","op","delete_file","err",coreerror.ErrPathTraversal)
		return coreerror.ErrPathTraversal
	}

	if err := os.Remove(target); err != nil {
		logger.Error(ctx,"failed delete file","layer","core.service","op","delete_file","err",err)
		return err
	}
	return nil

}