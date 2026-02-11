package app

import (
	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/create"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/update"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/delete"
	kelas_service_create "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	kelas_service_get "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
	kelas_service_update "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/update"
	kelas_service_delete "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/delete"
)

type KelasModule struct {
	GetService    *kelas_service_get.GetKelasService
	CreateService *kelas_service_create.CreateKelasService
	UpdateService *kelas_service_update.UpdateKelasService
	DeleteService *kelas_service_delete.DeleteKelasService

	GetHandler    *httpget.GetKelasHandler
	CreateHandler *httpcreate.CreateKelasHandler
	UpdateHandler *httpupdate.UpdateKelasHandler
	DeleteHandler *httpdelete.DeleteKelasHandler
}

func BuildKelasModule(infra *InfraModule, aktivitasUser *AktivitasUserModule) *KelasModule {
	svc := kelas_service_get.NewGetKelasService(infra.kelasRepo)
	createSvc := kelas_service_create.NewCreateKelasService(infra.kelasRepo)
	updateSvc := kelas_service_update.NewUpdateKelasService(infra.kelasRepo)
	deleteSvc := kelas_service_delete.NewDeleteKelasService(infra.kelasRepo)

	handler := httpget.NewGetKelasHandler(svc)
	createHandler := httpcreate.NewCreateKelasHandler(createSvc)
	updateHandler := httpupdate.NewUpdateKelasHandler(updateSvc, aktivitasUser.Service)
	deleteHandler := httpdelete.NewDeleteKelasHandler(deleteSvc, aktivitasUser.Service)

	return &KelasModule{
		GetService:    svc,
		GetHandler:    handler,
		CreateService: createSvc,
		DeleteService: deleteSvc,

		CreateHandler: createHandler,
		UpdateService: updateSvc,
		UpdateHandler: updateHandler,
		DeleteHandler: deleteHandler,
	}
}
