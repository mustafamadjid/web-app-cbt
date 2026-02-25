package app

import (
	delete_file_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/delete_file_system"
	delete_file_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/delete_file_system"
)

type DeleteFileModule struct {
	deleteFileRepo delete_file_repo.DeleteFileRepo
}

func BuildDeleteFileModule(uploadDir string, ) *DeleteFileModule {
	return &DeleteFileModule{
		deleteFileRepo : delete_file_service.NewDeleteFileService(uploadDir),
	}
}