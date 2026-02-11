package app

import (
	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/create"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/update"
	kelas_service_create "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	kelas_service_get "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
	kelas_service_update "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/update"
)

type KelasModule struct {
	GetService    *kelas_service_get.GetKelasService
	CreateService *kelas_service_create.CreateKelasService
	UpdateService *kelas_service_update.UpdateKelasService

	GetHandler    *httpget.GetKelasHandler
	CreateHandler *httpcreate.CreateKelasHandler
	UpdateHandler *httpupdate.UpdateKelasHandler
}

func BuildKelasModule(infra *InfraModule, aktivitasUser *AktivitasUserModule) *KelasModule {
	svc := kelas_service_get.NewGetKelasService(infra.kelasRepo)
	createSvc := kelas_service_create.NewCreateKelasService(infra.kelasRepo)
	updateSvc := kelas_service_update.NewUpdateKelasService(infra.kelasRepo)

	handler := httpget.NewGetKelasHandler(svc)
	createHandler := httpcreate.NewCreateKelasHandler(createSvc)
	updateHandler := httpupdate.NewUpdateKelasHandler(updateSvc, aktivitasUser.Service)

	return &KelasModule{
		GetService:    svc,
		GetHandler:    handler,
		CreateService: createSvc,
		CreateHandler: createHandler,
		UpdateService: updateSvc,
		UpdateHandler: updateHandler,
	}
}
