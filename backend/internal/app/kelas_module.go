package app

import (
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/get"

	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
)

type KelasModule struct {
	Service    *kelas_service.GetKelasService
	GetHandler *httpget.GetKelasHandler
}

func BuildKelasModule(infra *InfraModule) *KelasModule {
	svc := kelas_service.NewGetKelasService(infra.kelasRepo)
	handler := httpget.NewGetKelasHandler(svc)

	return &KelasModule{
		Service:    svc,
		GetHandler: handler,
	}
}
