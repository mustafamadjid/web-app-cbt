package app

import (
	httpresetpassword "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/reset_password"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/reset_password"
)

type ResetPasswordModule struct {
	Repo    outuser.UserResetPasswordRepo
	Service *user_service.ResetPasswordService
	Handler *httpresetpassword.ResetPasswordHandler
}

func BuildResetPasswordModule(infra *InfraModule, hasher out.PasswordHasher) *ResetPasswordModule {
	repo := infra.userResetPasswordRepo
	svc := user_service.NewResetPasswordService(repo, hasher)
	handler := httpresetpassword.NewResetPasswordHandler(svc)

	return &ResetPasswordModule{
		Repo:    repo,
		Service: svc,
		Handler: handler,
	}
}
