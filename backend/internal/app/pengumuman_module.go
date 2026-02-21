package app

import (
	"path/filepath"
	"strings"

	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/create"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/delete"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/update"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	pengumuman_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/create"
	pengumuman_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/delete"
	pengumuman_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/get"
	pengumuman_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/update"
)

type PengumumanModule struct {
	GetService    *pengumuman_get_service.GetPengumumanService
	CreateService *pengumuman_create_service.CreatePengumumanService
	UpdateService *pengumuman_update_service.UpdatePengumumanService
	DeleteService *pengumuman_delete_service.DeletePengumumanService

	GetHandler    *httpget.GetPengumumanHandler
	CreateHandler *httpcreate.CreatePengumumanHandler
	UpdateHandler *httpupdate.UpdatePengumumanHandler
	DeleteHandler *httpdelete.DeletePengumumanHandler
}

func BuildPengumumanModule(cfg Config, infra *InfraModule) *PengumumanModule {
	documentStore := httphelper.DocumentStore{
		Dir:      filepath.Join(cfg.DocumentStore.Dir, "pengumuman"),
		BaseURL:  cfg.DocumentStore.BaseURL,
		Route:    strings.TrimRight(cfg.DocumentStore.Route, "/") + "/pengumuman",
		MaxBytes: cfg.DocumentStore.MaxBytes,
	}

	getSvc := pengumuman_get_service.NewGetPengumumanService(infra.pengumumanRepo)
	createSvc := pengumuman_create_service.NewCreatePengumumanRepo(infra.pengumumanRepo)
	updateSvc := pengumuman_update_service.NewUpdatePengumumanService(infra.pengumumanRepo)
	deleteSvc := pengumuman_delete_service.NewDeletePengumumanService(infra.pengumumanRepo)

	getHandler := httpget.NewGetPengumumanHandler(getSvc)
	createHandler := httpcreate.NewCreatePengumumanHandler(createSvc, documentStore)
	updateHandler := httpupdate.NewUpdatePengumumanHandler(updateSvc, documentStore)
	deleteHandler := httpdelete.NewDeletePengumumanHandler(deleteSvc)

	return &PengumumanModule{
		GetService:    getSvc,
		CreateService: createSvc,
		UpdateService: updateSvc,
		DeleteService: deleteSvc,
		GetHandler:    getHandler,
		CreateHandler: createHandler,
		UpdateHandler: updateHandler,
		DeleteHandler: deleteHandler,
	}
}
