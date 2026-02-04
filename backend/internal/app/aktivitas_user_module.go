package app

import (
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/aktivitas_user/get"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
)

type AktivitasUserModule struct {
	Service    *aktivitas_user_service.AktivitasUserService
	GetHandler *httpget.AktivitasUserHandler
}

func BuildAktivitasUserModule(infra *InfraModule) *AktivitasUserModule {
	svc := aktivitas_user_service.NewAktivitasUserService(infra.aktivitasUser)
	handler := httpget.NewAktivitasUserHandler(svc)

	return &AktivitasUserModule{
		Service:    svc,
		GetHandler: handler,
	}
}
