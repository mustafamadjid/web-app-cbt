package app

import (
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/get"
	httplist "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list"
	ujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/get"
)

type UjianModule struct {
	GetService  *ujian_get_service.GetUjianService
	ListHandler *httplist.ListUjianHandler
	GetHandler  *httpget.GetUjianHandler
}

func BuildUjianModule(infra *InfraModule) *UjianModule {
	getSvc := ujian_get_service.NewGetujianService(infra.listUjianRepo)
	listHandler := httplist.NewListUjianHandler(getSvc)
	getHandler := httpget.NewGetUjianHandler(getSvc)

	return &UjianModule{
		GetService:  getSvc,
		ListHandler: listHandler,
		GetHandler:  getHandler,
	}
}
