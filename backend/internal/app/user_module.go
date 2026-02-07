package app

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"

	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/create"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/delete"
	httpget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/get"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/update"
	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"

	create "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	delete "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
	get "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
	update "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
)

type UserModule struct {
	CreateHandler   *httpcreate.UserHandler
	UpdateHandler   *httpupdate.UpdateHandler
	DeleteHandler   *httpdelete.DeleteHandler
	GetSiswaHandler *httpget.GetSiswaHandler
	GetGuruHandler  *httpget.GetGuruHandler
}

func BuildUserModule(cfg Config, infra *InfraModule, hasher out.PasswordHasher, aktivitasUser *AktivitasUserModule) *UserModule {
	store := httpx.ImageStore{
		Dir:      cfg.ImageStore.Dir,
		BaseURL:  cfg.ImageStore.BaseURL,
		Route:    cfg.ImageStore.Route,
		MaxBytes: cfg.ImageStore.MaxBytes,
	}

	createSvc := create.NewCreateGuruService(infra.Txm, hasher)
	updateSvc := update.NewUpdateGuruService(infra.Txm)
	deleteSvc := delete.NewDeleteUserService(infra.users)

	getSiswaSvc := get.NewGetListSiswaService(infra.profilSiswa)
	getGuruSvc := get.NewGetListGuruService(infra.profilGuru, infra.profilGuruRepo)

	handlerCreate := httpcreate.NewCreateUserHandler(createSvc, store, aktivitasUser.Service)
	handlerUpdate := httpupdate.NewUpdateUserHandler(updateSvc, store, aktivitasUser.Service)
	handlerDelete := httpdelete.NewDeleteUserHandler(deleteSvc, aktivitasUser.Service)

	handlerGetSiswa := httpget.NewGetSiswaHandler(getSiswaSvc)
	handlerGetGuru := httpget.NewGetGuruHandler(getGuruSvc)

	return &UserModule{
		CreateHandler:   handlerCreate,
		UpdateHandler:   handlerUpdate,
		DeleteHandler:   handlerDelete,
		GetSiswaHandler: handlerGetSiswa,
		GetGuruHandler:  handlerGetGuru,
	}
}
