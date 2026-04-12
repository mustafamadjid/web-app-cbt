package app

import (
	"time"

	httpactiveattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/active_attempt"
	httpcreateattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/create"
	httpexpireattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/expire"
	httpupdateattemptujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/update"
	httplistpesertasubmitted "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/list_peserta_submitted"
	httpgetjawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/jawaban_ujian/get_jawaban"
	httphasiljawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/jawaban_ujian/hasil_jawaban"
	httpsavejawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/jawaban_ujian/save_jawaban"
	httpkoreksiessay "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/koreksi_essay"
	httplist "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list"
	httplistsoal "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian"
	httplistsoalsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian_for_siswa"
	httplistessayungraded "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_ujian_essay_ungraded"
	httplistsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_ujian_siswa"
	httpstatistikujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/statistik_ujian"
	httpsubmitujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/submit_ujian"
	httpcreateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/create"
	httpdeleteujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/delete"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/get"
	httpupdateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/update"
	httpgetwaktuujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/waktu_ujian"
	activeattempt_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/active_attempt"
	attempt_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/create"
	attempt_list_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/list_peserta_submitted"
	attempt_submit_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/submit_ujian"
	attempt_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update"
	gradingujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading"
	essaygrading_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading/essay_grading"
	listessayungraded_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/list/list_ujian_essay_ungraded"
	gradingstatistikujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/statistik_ujian"
	gradingworker_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/worker"
	getjawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/get_jawaban"
	hasiljawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/hasil_jawaban"
	savejawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/save_jawaban"
	siswaujianlist_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list"
	siswaujiansoal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/soal_ujian"
	siswaujianwaktuselesai_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/waktu_selesai"
	ujian_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian"
	getstatistikujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/statistik_ujian"
	ujian_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/create"
	ujian_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/delete"
	ujian_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/get"
	ujian_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/update"
)

type UjianModule struct {
	CreateUjianService               *ujian_create_service.CreateUjianService
	AttemptUjianService              *attempt_create_service.AttemptUjianService
	GetActiveAttemptService          *activeattempt_service.GetActiveAttemptUjianService
	ListPesertaUjianSubmittedService *attempt_list_service.PesertaUjianSubmittedService
	UpdateAttemptUjianService        *attempt_update_service.SiswaUpdateAttemptUjianService
	ExpireAttemptUjianService        *attempt_update_service.ExpireAttemptUjianService
	SubmitUjianService               *attempt_submit_service.SubmitUjianService
	GradingUjianService              *gradingujian_service.GradingUjianService
	EssayGradingService              *essaygrading_service.EssayGradingUjianService
	ListUjianEssayUngradedService    *listessayungraded_service.ListUjianEssayUngradedService
	GetStatistikUjianService         *getstatistikujian_service.StatistikUjianService
	GradingWorker                    *gradingworker_service.GradingUjianWorkerRepo
	GetJawabanService                *getjawaban_service.SiswaGetJawabanUjianService
	HasilJawabanService              *hasiljawaban_service.HasilJawabanUjianService
	SaveJawabanService               *savejawaban_service.JawabanUjianService
	ListUjianSiswaService            *siswaujianlist_service.ListUjianSiswaService
	GetService                       *ujian_get_service.GetUjianService
	ListSoalUjianService             *ujian_soal_service.ListSoalUjianService
	ListSoalUjianSiswaService        *siswaujiansoal_service.ListSoalUjianSiswaService
	GetWaktuSelesaiService           *siswaujianwaktuselesai_service.GetWaktuSelesaiService
	UpdateUjianService               *ujian_update_service.UpdateUjianService
	DeleteUjianService               *ujian_delete_service.DeleteUjianService

	AttemptUjianHandler              *httpcreateattemptujian.AttemptUjianHandler
	GetActiveAttemptUjianHandler     *httpactiveattemptujian.GetActiveAttemptUjianHandler
	ListPesertaUjianSubmittedHandler *httplistpesertasubmitted.ListPesertaUjianSubmittedHandler
	UpdateAttemptUjianHandler        *httpupdateattemptujian.UpdateAttemptUjianHandler
	ExpireAttemptUjianHandler        *httpexpireattemptujian.ExpireAttemptUjianHandler
	SubmitUjianHandler               *httpsubmitujian.SubmitUjianHandler
	GetJawabanUjianHandler           *httpgetjawaban.GetJawabanUjianHandler
	HasilJawabanUjianHandler         *httphasiljawaban.HasilJawabanUjianHandler
	SaveJawabanUjianHandler          *httpsavejawaban.SaveJawabanUjianHandler
	KoreksiEssayHandler              *httpkoreksiessay.KoreksiEssayHandler
	ListUjianEssayUngradedHandler    *httplistessayungraded.ListUjianEssayUngradedHandler
	GetStatistikUjianHandler         *httpstatistikujian.GetStatistikUjianHandler
	CreateUjianHandler               *httpcreateujian.CreateRuangUjianHandler
	ListHandler                      *httplist.ListUjianHandler
	ListUjianSiswaHandler            *httplistsiswa.ListUjianSiswaHandler
	ListSoalUjianHandler             *httplistsoal.ListSoalUjianHandler
	ListSoalUjianSiswaHandler        *httplistsoalsiswa.ListSoalUjianSiswaHandler
	GetWaktuSelesaiUjianHandler      *httpgetwaktuujian.GetWaktuSelesaiUjianHandler
	GetHandler                       *httpget.GetUjianHandler
	UpdateUjianHandler               *httpupdateujian.UpdateUjianHandler
	DeleteUjianHandler               *httpdeleteujian.DeleteUjianHandler
}

func BuildUjianModule(infra *InfraModule) *UjianModule {
	createUjianSvc := ujian_create_service.NewCreateUjianService(infra.ujianRepo)
	attemptUjianSvc := attempt_create_service.NewAttemptUjianService(
		infra.siswaUjianChecker,
		infra.attemptUjianRepo,
	)
	getActiveAttemptSvc := activeattempt_service.NewGetActiveAttemptUjianService(infra.ujianSiswaRepo)
	listPesertaUjianSubmittedSvc := attempt_list_service.NewPesertaUjianSubmittedService(infra.attemptUjianRepo)
	getJawabanSvc := getjawaban_service.NewGetJawabanUjianService(infra.jawabanUjianRepo)
	siswaGetJawabanSvc := getjawaban_service.NewSiswaGetJawabanUjianService(infra.siswaUjianChecker, getJawabanSvc)
	hasilJawabanSvc := hasiljawaban_service.NewHasilJawabanUjianService(infra.jawabanUjianRepo)
	saveJawabanSvc := savejawaban_service.NewJawabanUjianService(infra.jawabanUjianRepo)
	listUjianSiswaSvc := siswaujianlist_service.NewListUjianSiswaService(infra.ujianSiswaRepo)
	getSvc := ujian_get_service.NewGetujianService(infra.listUjianRepo)
	listSoalUjianSvc := ujian_soal_service.NewListSoalUjianService(infra.soalUjianRepo)
	listSoalUjianSiswaSvc := siswaujiansoal_service.NewListSoalUjianSiswaService(infra.soalUjianRepo)
	getWaktuSelesaiSvc := siswaujianwaktuselesai_service.NewGetWaktuSelesaiService(infra.ujianSiswaRepo)
	updateAttemptUjianSvc := attempt_update_service.NewUpdateAttemptUjianService(infra.attemptUjianRepo)
	siswaUpdateAttemptUjianSvc := attempt_update_service.NewSiswaUpdateAttemptUjianService(infra.siswaUjianChecker, updateAttemptUjianSvc)
	expireAttemptUjianSvc := attempt_update_service.NewExpireAttemptUjianService(updateAttemptUjianSvc)
	submitUjianSvc := attempt_submit_service.NewSubmitUjianService(infra.attemptUjianRepo, infra.siswaUjianChecker)
	gradingUjianSvc := gradingujian_service.NewGradingUjianService(
		infra.jawabanUjianRepo,
		infra.soalUjianRepo,
		infra.bankSoalRepo,
		infra.ujianRepo,
		infra.gradingRepo,
	)
	essayGradingSvc := essaygrading_service.NewEssayGradingUjianService(infra.gradingRepo)
	listEssayUngradedSvc := listessayungraded_service.NewListUjianEssayUngradedService(infra.listGradingRepo)
	getStatistikUjianSvc := getstatistikujian_service.NewStatistikUjianService(infra.statistikUjianRepo)
	statistikUjianSvc := gradingstatistikujian_service.NewStatistikUjianService(infra.gradingRepo)
	gradingWorker := gradingworker_service.NewGradingUjianWorkerService(
		infra.gradingWorkerRepo,
		gradingworker_service.NewCompositeGradingUjianExecutor(gradingUjianSvc, statistikUjianSvc),
		5*time.Second,
	)
	updateUjianSvc := ujian_update_service.NewUpdateUjianService(infra.ujianRepo)
	deleteUjianSvc := ujian_delete_service.NewDeleteUjianService(infra.ujianRepo)

	attemptUjianHandler := httpcreateattemptujian.NewAttemptUjianHandler(attemptUjianSvc)
	getActiveAttemptHandler := httpactiveattemptujian.NewGetActiveAttemptUjianHandler(getActiveAttemptSvc)
	listPesertaUjianSubmittedHandler := httplistpesertasubmitted.NewListPesertaUjianSubmittedHandler(listPesertaUjianSubmittedSvc)
	updateAttemptUjianHandler := httpupdateattemptujian.NewUpdateAttemptUjianHandler(siswaUpdateAttemptUjianSvc)
	expireAttemptUjianHandler := httpexpireattemptujian.NewExpireAttemptUjianHandler(expireAttemptUjianSvc)
	submitUjianHandler := httpsubmitujian.NewSubmitUjianHandler(submitUjianSvc)
	getJawabanHandler := httpgetjawaban.NewGetJawabanUjianHandler(siswaGetJawabanSvc)
	hasilJawabanHandler := httphasiljawaban.NewHasilJawabanUjianHandler(hasilJawabanSvc)
	saveJawabanHandler := httpsavejawaban.NewSaveJawabanUjianHandler(saveJawabanSvc)
	koreksiEssayHandler := httpkoreksiessay.NewKoreksiEssayHandler(essayGradingSvc)
	listEssayUngradedHandler := httplistessayungraded.NewListUjianEssayUngradedHandler(listEssayUngradedSvc)
	getStatistikUjianHandler := httpstatistikujian.NewGetStatistikUjianHandler(getStatistikUjianSvc)
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
		AttemptUjianService:              attemptUjianSvc,
		GetActiveAttemptService:          getActiveAttemptSvc,
		ListPesertaUjianSubmittedService: listPesertaUjianSubmittedSvc,
		GetJawabanService:                siswaGetJawabanSvc,
		HasilJawabanService:              hasilJawabanSvc,
		SaveJawabanService:               saveJawabanSvc,
		CreateUjianService:               createUjianSvc,
		ListUjianSiswaService:            listUjianSiswaSvc,
		GetService:                       getSvc,
		UpdateAttemptUjianService:        siswaUpdateAttemptUjianSvc,
		ExpireAttemptUjianService:        expireAttemptUjianSvc,
		SubmitUjianService:               submitUjianSvc,
		GradingUjianService:              gradingUjianSvc,
		EssayGradingService:              essayGradingSvc,
		ListUjianEssayUngradedService:    listEssayUngradedSvc,
		GetStatistikUjianService:         getStatistikUjianSvc,
		GradingWorker:                    gradingWorker,
		ListSoalUjianService:             listSoalUjianSvc,
		ListSoalUjianSiswaService:        listSoalUjianSiswaSvc,
		GetWaktuSelesaiService:           getWaktuSelesaiSvc,
		UpdateUjianService:               updateUjianSvc,
		DeleteUjianService:               deleteUjianSvc,
		AttemptUjianHandler:              attemptUjianHandler,
		GetActiveAttemptUjianHandler:     getActiveAttemptHandler,
		ListPesertaUjianSubmittedHandler: listPesertaUjianSubmittedHandler,
		UpdateAttemptUjianHandler:        updateAttemptUjianHandler,
		ExpireAttemptUjianHandler:        expireAttemptUjianHandler,
		SubmitUjianHandler:               submitUjianHandler,
		GetJawabanUjianHandler:           getJawabanHandler,
		HasilJawabanUjianHandler:         hasilJawabanHandler,
		SaveJawabanUjianHandler:          saveJawabanHandler,
		KoreksiEssayHandler:              koreksiEssayHandler,
		ListUjianEssayUngradedHandler:    listEssayUngradedHandler,
		GetStatistikUjianHandler:         getStatistikUjianHandler,
		CreateUjianHandler:               createUjianHandler,
		ListHandler:                      listHandler,
		ListUjianSiswaHandler:            listUjianSiswaHandler,
		ListSoalUjianHandler:             listSoalUjianHandler,
		ListSoalUjianSiswaHandler:        listSoalUjianSiswaHandler,
		GetWaktuSelesaiUjianHandler:      getWaktuSelesaiHandler,
		GetHandler:                       getHandler,
		UpdateUjianHandler:               updateUjianHandler,
		DeleteUjianHandler:               deleteUjianHandler,
	}
}
