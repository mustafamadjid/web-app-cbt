package app

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/create"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/delete"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/update"

	create "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	delete "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
	update "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
)

type UserModule struct {
	CreateHandler *httpcreate.UserHandler
	UpdateHandler *httpupdate.UpdateHandler	
	DeleteHandler *httpdelete.DeleteHandler
}

func BuildUserModule(cfg Config,infra *InfraModule, hasher out.PasswordHasher) *UserModule {
	store := httpx.ImageStore{
		Dir: cfg.ImageStore.Dir,
		BaseURL: cfg.ImageStore.BaseURL,
		Route: cfg.ImageStore.Route,
		MaxBytes: cfg.ImageStore.MaxBytes,
	}

	createSvc := create.NewCreateGuruService(infra.Txm, hasher)
	updateSvc := update.NewUpdateGuruService(infra.Txm)
	deleteSvc := delete.NewDeleteUserService(infra.users)

	handlerCreate := httpcreate.NewCreateUserHandler(createSvc, store)
	handlerUpdate := httpupdate.NewUpdateUserHandler(updateSvc, store)
	handlerDelete := httpdelete.NewDeleteUserHandler(deleteSvc)

	return &UserModule{CreateHandler: handlerCreate, UpdateHandler: handlerUpdate, DeleteHandler: handlerDelete}
}