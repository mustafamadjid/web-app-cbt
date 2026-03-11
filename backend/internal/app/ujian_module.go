package app

import (
	httpcreateattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/create/attempt_ujian"
	httpcreateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/create/ujian"
	httpdeleteujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/delete/ujian"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/get"
	httplist "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list"
	httplistsoal "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian"
	httpupdateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/update/ujian"
	ujian_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/create"
	ujian_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/delete"
	ujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/get"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian"
	ujian_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian"
	ujian_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/update"
)

type UjianModule struct {
	CreateUjianService   *ujian_create_service.CreateUjianService
	AttemptUjianService  *siswaujian_service.AttemptUjianService
	GetService           *ujian_get_service.GetUjianService
	ListSoalUjianService *ujian_soal_service.ListSoalUjianService
	UpdateUjianService   *ujian_update_service.UpdateUjianService
	DeleteUjianService   *ujian_delete_service.DeleteUjianService

	AttemptUjianHandler  *httpcreateattemptujian.AttemptUjianHandler
	CreateUjianHandler   *httpcreateujian.CreateRuangUjianHandler
	ListHandler          *httplist.ListUjianHandler
	ListSoalUjianHandler *httplistsoal.ListSoalUjianHandler
	GetHandler           *httpget.GetUjianHandler
	UpdateUjianHandler   *httpupdateujian.UpdateUjianHandler
	DeleteUjianHandler   *httpdeleteujian.DeleteUjianHandler
}

func BuildUjianModule(infra *InfraModule) *UjianModule {
	createUjianSvc := ujian_create_service.NewCreateUjianService(
		ujian_create_service.NewUjianRepository(infra.ujianRepo),
	)
	attemptUjianSvc := siswaujian_service.NewAttemptUjianService(
		infra.siswaUjianChecker,
		infra.attemptUjianRepo,
	)
	getSvc := ujian_get_service.NewGetujianService(infra.listUjianRepo)
	listSoalUjianSvc := ujian_soal_service.NewListSoalUjianService(infra.soalUjianRepo)
	updateUjianSvc := ujian_update_service.NewUpdateUjianService(
		ujian_update_service.NewUjianRepository(infra.ujianRepo),
	)
	deleteUjianSvc := ujian_delete_service.NewDeleteUjianService(
		ujian_delete_service.NewUjianRepository(infra.ujianRepo),
	)

	attemptUjianHandler := httpcreateattemptujian.NewAttemptUjianHandler(attemptUjianSvc)
	createUjianHandler := httpcreateujian.NewCreateUjianHandler(createUjianSvc)
	listHandler := httplist.NewListUjianHandler(getSvc)
	listSoalUjianHandler := httplistsoal.NewListSoalUjianHandler(listSoalUjianSvc)
	getHandler := httpget.NewGetUjianHandler(getSvc)
	updateUjianHandler := httpupdateujian.NewUpdateUjianHandler(updateUjianSvc)
	deleteUjianHandler := httpdeleteujian.NewDeleteUjianHandler(deleteUjianSvc)

	return &UjianModule{
		AttemptUjianService:  attemptUjianSvc,
		CreateUjianService:   createUjianSvc,
		GetService:           getSvc,
		ListSoalUjianService: listSoalUjianSvc,
		UpdateUjianService:   updateUjianSvc,
		DeleteUjianService:   deleteUjianSvc,
		AttemptUjianHandler:  attemptUjianHandler,
		CreateUjianHandler:   createUjianHandler,
		ListHandler:          listHandler,
		ListSoalUjianHandler: listSoalUjianHandler,
		GetHandler:           getHandler,
		UpdateUjianHandler:   updateUjianHandler,
		DeleteUjianHandler:   deleteUjianHandler,
	}
}
