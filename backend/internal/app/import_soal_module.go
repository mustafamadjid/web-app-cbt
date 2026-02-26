package app

import (
	"path/filepath"
	"time"

	importhandler "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/import"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/create_job"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/get_job"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/worker"
)

type ImportSoalModule struct {
	CreateJobService *create_job.CreateJobService
	GetJobService    *get_job.GetJobService
	ImportHandler    *importhandler.ImportHandler
	GetJobHandler    *importhandler.GetJobHandler
	Worker           *worker.Worker
}

func BuildImportSoalModule(infra *InfraModule, cfg Config, logger corelog.Logger) *ImportSoalModule {
	createJobSvc := create_job.NewCreateJobService(infra.importSoalJobRepo)
	getJobSvc := get_job.NewGetJobService(infra.importSoalJobRepo)

	docxUploadDir := filepath.Join(cfg.DocumentStore.Dir, "import_soal")
	importHandler := importhandler.NewImportHandler(createJobSvc, docxUploadDir)
	getJobHandler := importhandler.NewGetJobHandler(getJobSvc)

	imageDir := filepath.Join(cfg.ImageStore.Dir, "bank_soal")
	imageRoute := cfg.ImageStore.Route + "/bank_soal"

	w := worker.NewWorker(
		infra.importSoalJobRepo,
		infra.isiSoalBatchRepo,
		imageDir,
		imageRoute,
		5*time.Second,
		logger,
	)

	return &ImportSoalModule{
		CreateJobService: createJobSvc,
		GetJobService:    getJobSvc,
		ImportHandler:    importHandler,
		GetJobHandler:    getJobHandler,
		Worker:           w,
	}
}
