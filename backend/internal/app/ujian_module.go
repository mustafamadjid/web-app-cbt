package app

import (
	httpactiveattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/active_attempt"
	httpcreateattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/create"
	httpexpireattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/expire"
	httpupdateattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/update"
	httpsavejawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/jawaban_ujian/save_jawaban"
	httplist "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list"
	httplistsoal "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian"
	httplistsoalsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian_for_siswa"
	httplistsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_ujian_siswa"
	httpcreateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/create"
	httpdeleteujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/delete"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/get"
	httpupdateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/update"
	httpgetwaktuujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/waktu_ujian"
	attempt_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian"
	activeattempt_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/active_attempt"
	savejawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/save_jawaban"
	siswaujianlist_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list"
	siswaujiansoal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/soal_ujian"
	siswaujianwaktuselesai_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/waktu_selesai"
	ujian_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian"
	ujian_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/create"
	ujian_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/delete"
	ujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/get"
	ujian_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/update"
)

type UjianModule struct {
	CreateUjianService        *ujian_create_service.CreateUjianService
	AttemptUjianService       *siswaujian_service.AttemptUjianService
	GetActiveAttemptService   *activeattempt_service.GetActiveAttemptUjianService
	UpdateAttemptUjianService *attempt_update_service.SiswaUpdateAttemptUjianService
	ExpireAttemptUjianService *attempt_update_service.ExpireAttemptUjianService
	SaveJawabanService        *savejawaban_service.JawabanUjianService
	ListUjianSiswaService     *siswaujianlist_service.ListUjianSiswaService
	GetService                *ujian_get_service.GetUjianService
	ListSoalUjianService      *ujian_soal_service.ListSoalUjianService
	ListSoalUjianSiswaService *siswaujiansoal_service.ListSoalUjianSiswaService
	GetWaktuSelesaiService    *siswaujianwaktuselesai_service.GetWaktuSelesaiService
	UpdateUjianService        *ujian_update_service.UpdateUjianService
	DeleteUjianService        *ujian_delete_service.DeleteUjianService

	AttemptUjianHandler          *httpcreateattemptujian.AttemptUjianHandler
	GetActiveAttemptUjianHandler *httpactiveattemptujian.GetActiveAttemptUjianHandler
	UpdateAttemptUjianHandler    *httpupdateattemptujian.UpdateAttemptUjianHandler
	ExpireAttemptUjianHandler    *httpexpireattemptujian.ExpireAttemptUjianHandler
	SaveJawabanUjianHandler      *httpsavejawaban.SaveJawabanUjianHandler
	CreateUjianHandler           *httpcreateujian.CreateRuangUjianHandler
	ListHandler                  *httplist.ListUjianHandler
	ListUjianSiswaHandler        *httplistsiswa.ListUjianSiswaHandler
	ListSoalUjianHandler         *httplistsoal.ListSoalUjianHandler
	ListSoalUjianSiswaHandler    *httplistsoalsiswa.ListSoalUjianSiswaHandler
	GetWaktuSelesaiUjianHandler  *httpgetwaktuujian.GetWaktuSelesaiUjianHandler
	GetHandler                   *httpget.GetUjianHandler
	UpdateUjianHandler           *httpupdateujian.UpdateUjianHandler
	DeleteUjianHandler           *httpdeleteujian.DeleteUjianHandler
}

func BuildUjianModule(infra *InfraModule) *UjianModule {
	createUjianSvc := ujian_create_service.NewCreateUjianService(infra.ujianRepo)
	attemptUjianSvc := siswaujian_service.NewAttemptUjianService(
		infra.siswaUjianChecker,
		infra.attemptUjianRepo,
	)
	getActiveAttemptSvc := activeattempt_service.NewGetActiveAttemptUjianService(infra.activeAttemptUjianRepo)
	saveJawabanSvc := savejawaban_service.NewJawabanUjianService(infra.jawabanUjianRepo)
	listUjianSiswaSvc := siswaujianlist_service.NewListUjianSiswaService(infra.listUjianSiswaRepo)
	getSvc := ujian_get_service.NewGetujianService(infra.listUjianRepo)
	listSoalUjianSvc := ujian_soal_service.NewListSoalUjianService(infra.soalUjianRepo)
	listSoalUjianSiswaSvc := siswaujiansoal_service.NewListSoalUjianSiswaService(infra.soalUjianRepo)
	getWaktuSelesaiSvc := siswaujianwaktuselesai_service.NewGetWaktuSelesaiService(infra.waktuSelesaiUjianRepo)
	updateAttemptUjianSvc := attempt_update_service.NewUpdateAttemptUjianService(infra.attemptUjianRepo)
	siswaUpdateAttemptUjianSvc := attempt_update_service.NewSiswaUpdateAttemptUjianService(infra.siswaUjianChecker, updateAttemptUjianSvc)
	expireAttemptUjianSvc := attempt_update_service.NewExpireAttemptUjianService(updateAttemptUjianSvc)
	updateUjianSvc := ujian_update_service.NewUpdateUjianService(infra.ujianRepo)
	deleteUjianSvc := ujian_delete_service.NewDeleteUjianService(infra.ujianRepo)

	attemptUjianHandler := httpcreateattemptujian.NewAttemptUjianHandler(attemptUjianSvc)
	getActiveAttemptHandler := httpactiveattemptujian.NewGetActiveAttemptUjianHandler(getActiveAttemptSvc)
	updateAttemptUjianHandler := httpupdateattemptujian.NewUpdateAttemptUjianHandler(siswaUpdateAttemptUjianSvc)
	expireAttemptUjianHandler := httpexpireattemptujian.NewExpireAttemptUjianHandler(expireAttemptUjianSvc)
	saveJawabanHandler := httpsavejawaban.NewSaveJawabanUjianHandler(saveJawabanSvc)
	createUjianHandler := httpcreateujian.NewCreateUjianHandler(createUjianSvc)
	listHandler := httplist.NewListUjianHandler(getSvc)
	listUjianSiswaHandler := httplistsiswa.NewListUjianSiswaHandler(listUjianSiswaSvc)
	listSoalUjianHandler := httplistsoal.NewListSoalUjianHandler(listSoalUjianSvc)
	listSoalUjianSiswaHandler := httplistsoalsiswa.NewListSoalUjianSiswaHandler(listSoalUjianSiswaSvc)
	getWaktuSelesaiHandler := httpgetwaktuujian.NewGetWaktuSelesaiUjianHandler(getWaktuSelesaiSvc)
	getHandler := httpget.NewGetUjianHandler(getSvc)
	updateUjianHandler := httpupdateujian.NewUpdateUjianHandler(updateUjianSvc)
	deleteUjianHandler := httpdeleteujian.NewDeleteUjianHandler(deleteUjianSvc)

	return &UjianModule{
		AttemptUjianService:          attemptUjianSvc,
		GetActiveAttemptService:      getActiveAttemptSvc,
		SaveJawabanService:           saveJawabanSvc,
		CreateUjianService:           createUjianSvc,
		ListUjianSiswaService:        listUjianSiswaSvc,
		GetService:                   getSvc,
		UpdateAttemptUjianService:    siswaUpdateAttemptUjianSvc,
		ExpireAttemptUjianService:    expireAttemptUjianSvc,
		ListSoalUjianService:         listSoalUjianSvc,
		ListSoalUjianSiswaService:    listSoalUjianSiswaSvc,
		GetWaktuSelesaiService:       getWaktuSelesaiSvc,
		UpdateUjianService:           updateUjianSvc,
		DeleteUjianService:           deleteUjianSvc,
		AttemptUjianHandler:          attemptUjianHandler,
		GetActiveAttemptUjianHandler: getActiveAttemptHandler,
		UpdateAttemptUjianHandler:    updateAttemptUjianHandler,
		ExpireAttemptUjianHandler:    expireAttemptUjianHandler,
		SaveJawabanUjianHandler:      saveJawabanHandler,
		CreateUjianHandler:           createUjianHandler,
		ListHandler:                  listHandler,
		ListUjianSiswaHandler:        listUjianSiswaHandler,
		ListSoalUjianHandler:         listSoalUjianHandler,
		ListSoalUjianSiswaHandler:    listSoalUjianSiswaHandler,
		GetWaktuSelesaiUjianHandler:  getWaktuSelesaiHandler,
		GetHandler:                   getHandler,
		UpdateUjianHandler:           updateUjianHandler,
		DeleteUjianHandler:           deleteUjianHandler,
	}
}
