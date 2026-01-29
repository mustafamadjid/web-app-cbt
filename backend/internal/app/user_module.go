package app

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"

	httpcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/create"
	httpupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/update"
	httpdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/delete"
	
	create "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	update "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
	delete "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
)

type UserModule struct {
	CreateHandler *httpcreate.UserHandler
	UpdateHandler *httpupdate.UpdateHandler	
	DeleteHandler *httpdelete.DeleteHandler
}

func BuildUserModule(infra *InfraModule, hasher out.PasswordHasher) *UserModule {
	createSvc := create.NewCreateGuruService(infra.Txm, hasher)
	updateSvc := update.NewUpdateGuruService(infra.Txm)
	deleteSvc := delete.NewDeleteUserService(infra.users)

	handlerCreate := httpcreate.NewCreateUserHandler(createSvc)
	handlerUpdate := httpupdate.NewUpdateUserHandler(updateSvc)
	handlerDelete := httpdelete.NewDeleteUserHandler(deleteSvc)

	return &UserModule{CreateHandler: handlerCreate, UpdateHandler: handlerUpdate, DeleteHandler: handlerDelete}
}