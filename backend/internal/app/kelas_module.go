package app

import (
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/get"
	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/create"

	kelas_service_get "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
	kelas_service_create "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	
)

type KelasModule struct {
	GetService    *kelas_service_get.GetKelasService
	CreateService *kelas_service_create.CreateKelasService

	GetHandler *httpget.GetKelasHandler
	CreateHandler *httpcreate.CreateKelasHandler
}

func BuildKelasModule(infra *InfraModule) *KelasModule {
	svc := kelas_service_get.NewGetKelasService(infra.kelasRepo)
	createSvc := kelas_service_create.NewCreateKelasService(infra.kelasRepo)

	handler := httpget.NewGetKelasHandler(svc)
	createHandler := httpcreate.NewCreateKelasHandler(createSvc)

	return &KelasModule{
		GetService:    svc,
		GetHandler: handler,
		CreateService:createSvc,
		CreateHandler: createHandler,
	}
}
