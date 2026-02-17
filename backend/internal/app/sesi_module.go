package app

import (
	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/create"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/delete"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/update"
	sesi_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
	sesi_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/delete"
	sesi_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/get"
	sesi_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/update"
)

type SesiModule struct {
	GetService    *sesi_get_service.GetSesiService
	CreateService *sesi_create_service.CreateSesiService
	UpdateService *sesi_update_service.UpdateSesiService
	DeleteService *sesi_delete_service.DeleteSesiService

	GetHandler    *httpget.GetSesiHandler
	CreateHandler *httpcreate.CreateSesiHandler
	UpdateHandler *httpupdate.UpdateSesiHandler
	DeleteHandler *httpdelete.DeleteSesiHandler
}

func BuildSesiModule(infra *InfraModule) *SesiModule {
	getSvc := sesi_get_service.NewGetSesiService(infra.sesiRepo)
	createSvc := sesi_create_service.NewCreateSesiService(infra.sesiRepo)
	updateSvc := sesi_update_service.NewUpdateSesiService(infra.sesiRepo)
	deleteSvc := sesi_delete_service.NewDeleteSesiService(infra.sesiRepo)

	getHandler := httpget.NewGetSesiHandler(getSvc)
	createHandler := httpcreate.NewCreateSesiHandler(createSvc)
	updateHandler := httpupdate.NewUpdateSesiHandler(updateSvc)
	deleteHandler := httpdelete.NewDeleteSesiHandler(deleteSvc)

	return &SesiModule{
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
