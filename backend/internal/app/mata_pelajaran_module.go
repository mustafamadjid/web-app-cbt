package app

import (
	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/create"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/delete"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/update"
	mapel_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/create"
	mapel_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/delete"
	mapel_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/get"
	mapel_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/update"
)

type MataPelajaranModule struct {
	GetService    *mapel_get_service.GetMapelRepo
	CreateService *mapel_create_service.CreateMapelRepo
	UpdateService *mapel_update_service.UpdateMapelRepo
	DeleteService *mapel_delete_service.DeleteMapelRepo

	GetHandler    *httpget.GetMapelHandler
	CreateHandler *httpcreate.CreateMapelHandler
	UpdateHandler *httpupdate.UpdateMapelHandler
	DeleteHandler *httpdelete.DeleteMapelHandler
}

func BuildMataPelajaranModule(infra *InfraModule) *MataPelajaranModule {
	getSvc := mapel_get_service.NewGetMapelService(infra.mapelRepo)
	createSvc := mapel_create_service.NewMapelService(infra.mapelRepo)
	updateSvc := mapel_update_service.NewUpdateMapelService(infra.mapelRepo)
	deleteSvc := mapel_delete_service.NewDeleteMapelService(infra.mapelRepo)

	getHandler := httpget.NewGetMapelHandler(getSvc)
	createHandler := httpcreate.NewCreateMapelHandler(createSvc)
	updateHandler := httpupdate.NewUpdateMapelHandler(updateSvc)
	deleteHandler := httpdelete.NewDeleteMapelHandler(deleteSvc)

	return &MataPelajaranModule{
		GetService:    getSvc,
		CreateService: createSvc,
		UpdateService: updateSvc,
		DeleteService: deleteSvc,
		GetHandler:    getHandler,
		CreateHandler: createHandler,
		UpdateHandler: updateHandler,
		DeleteHandler: deleteHandler,
	}
}
