package app

import (
	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/create"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/delete"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/update"
	ruangujian_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/create"
	ruangujian_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/delete"
	ruangujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/get"
	ruangujian_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/update"
)

type RuangUjianModule struct {
	GetService    *ruangujian_get_service.GetRuangUjianRepo
	CreateService *ruangujian_create_service.CreateRuangUjianService
	UpdateService *ruangujian_update_service.UpdateRuangUjianService
	DeleteService *ruangujian_delete_service.DeleteRuangUjianService

	GetHandler    *httpget.GetRuangUjianHandler
	CreateHandler *httpcreate.CreateRuangUjianHandler
	UpdateHandler *httpupdate.UpdateRuangUjianHandler
	DeleteHandler *httpdelete.DeleteRuangUjianHandler
}

func BuildRuangUjianModule(infra *InfraModule) *RuangUjianModule {
	getSvc := ruangujian_get_service.NewGetRuangUjianService(infra.ruangUjianRepo)
	createSvc := ruangujian_create_service.NewRuangUjianService(infra.ruangUjianRepo)
	updateSvc := ruangujian_update_service.NewUpdateRuangUjianService(infra.ruangUjianRepo)
	deleteSvc := ruangujian_delete_service.NewDeleteRuangUjianService(infra.ruangUjianRepo)

	getHandler := httpget.NewGetRuangUjianHandler(getSvc)
	createHandler := httpcreate.NewCreateRuangUjianHandler(createSvc)
	updateHandler := httpupdate.NewUpdateRuangUjianHandler(updateSvc)
	deleteHandler := httpdelete.NewDeleteRuangUjianHandler(deleteSvc)

	return &RuangUjianModule{
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
