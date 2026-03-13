package app

import (
	httpcreateattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/create/attempt_ujian"
	httpcreateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/create/ujian"
	httpdeleteujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/delete/ujian"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/get"
	httplist "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list"
	httplistsoal "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian"
	httplistsoalsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian_for_siswa"
	httplistsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_ujian_siswa"
	httpupdateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/update/ujian"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian"
	siswaujianlist_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list"
	ujian_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian"
	ujian_soal_siswa_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian/siswa"
	ujian_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/create"
	ujian_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/delete"
	ujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/get"
	ujian_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/update"
)

type UjianModule struct {
	CreateUjianService        *ujian_create_service.CreateUjianService
	AttemptUjianService       *siswaujian_service.AttemptUjianService
	ListUjianSiswaService     *siswaujianlist_service.ListUjianSiswaService
	GetService                *ujian_get_service.GetUjianService
	ListSoalUjianService      *ujian_soal_service.ListSoalUjianService
	ListSoalUjianSiswaService *ujian_soal_siswa_service.ListSoalUjianSiswaService
	UpdateUjianService        *ujian_update_service.UpdateUjianService
	DeleteUjianService        *ujian_delete_service.DeleteUjianService

	AttemptUjianHandler       *httpcreateattemptujian.AttemptUjianHandler
	CreateUjianHandler        *httpcreateujian.CreateRuangUjianHandler
	ListHandler               *httplist.ListUjianHandler
	ListUjianSiswaHandler     *httplistsiswa.ListUjianSiswaHandler
	ListSoalUjianHandler      *httplistsoal.ListSoalUjianHandler
	ListSoalUjianSiswaHandler *httplistsoalsiswa.ListSoalUjianSiswaHandler
	GetHandler                *httpget.GetUjianHandler
	UpdateUjianHandler        *httpupdateujian.UpdateUjianHandler
	DeleteUjianHandler        *httpdeleteujian.DeleteUjianHandler
}

func BuildUjianModule(infra *InfraModule) *UjianModule {
	createUjianSvc := ujian_create_service.NewCreateUjianService(infra.ujianRepo)
	attemptUjianSvc := siswaujian_service.NewAttemptUjianService(
		infra.siswaUjianChecker,
		infra.attemptUjianRepo,
	)
	listUjianSiswaSvc := siswaujianlist_service.NewListUjianSiswaService(infra.listUjianSiswaRepo)
	getSvc := ujian_get_service.NewGetujianService(infra.listUjianRepo)
	listSoalUjianSvc := ujian_soal_service.NewListSoalUjianService(infra.soalUjianRepo)
	listSoalUjianSiswaSvc := ujian_soal_siswa_service.NewListSoalUjianSiswaService(infra.soalUjianRepo)
	updateUjianSvc := ujian_update_service.NewUpdateUjianService(infra.ujianRepo)
	deleteUjianSvc := ujian_delete_service.NewDeleteUjianService(infra.ujianRepo)

	attemptUjianHandler := httpcreateattemptujian.NewAttemptUjianHandler(attemptUjianSvc)
	createUjianHandler := httpcreateujian.NewCreateUjianHandler(createUjianSvc)
	listHandler := httplist.NewListUjianHandler(getSvc)
	listUjianSiswaHandler := httplistsiswa.NewListUjianSiswaHandler(listUjianSiswaSvc)
	listSoalUjianHandler := httplistsoal.NewListSoalUjianHandler(listSoalUjianSvc)
	listSoalUjianSiswaHandler := httplistsoalsiswa.NewListSoalUjianSiswaHandler(listSoalUjianSiswaSvc)
	getHandler := httpget.NewGetUjianHandler(getSvc)
	updateUjianHandler := httpupdateujian.NewUpdateUjianHandler(updateUjianSvc)
	deleteUjianHandler := httpdeleteujian.NewDeleteUjianHandler(deleteUjianSvc)

	return &UjianModule{
		AttemptUjianService:       attemptUjianSvc,
		CreateUjianService:        createUjianSvc,
		ListUjianSiswaService:     listUjianSiswaSvc,
		GetService:                getSvc,
		ListSoalUjianService:      listSoalUjianSvc,
		ListSoalUjianSiswaService: listSoalUjianSiswaSvc,
		UpdateUjianService:        updateUjianSvc,
		DeleteUjianService:        deleteUjianSvc,
		AttemptUjianHandler:       attemptUjianHandler,
		CreateUjianHandler:        createUjianHandler,
		ListHandler:               listHandler,
		ListUjianSiswaHandler:     listUjianSiswaHandler,
		ListSoalUjianHandler:      listSoalUjianHandler,
		ListSoalUjianSiswaHandler: listSoalUjianSiswaHandler,
		GetHandler:                getHandler,
		UpdateUjianHandler:        updateUjianHandler,
		DeleteUjianHandler:        deleteUjianHandler,
	}
}
