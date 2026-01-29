package out

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	createSvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	updateSvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
)

type CreateUserService interface {
	CreateGuru(ctx context.Context, cmd createSvc.CreateGuruCmd, actor user.Actor) (createSvc.CreateGuruRes, error)
	CreateSiswa(ctx context.Context, cmd createSvc.CreateSiswaCmd, actor user.Actor) (createSvc.CreateSiswaRes, error)
}

type UpdateUserService interface {
	UpdateGuru(ctx context.Context, cmd updateSvc.UpdateGuruCmd, actor user.Actor) error
	UpdateSiswa(ctx context.Context, cmd updateSvc.UpdateSiswaCmd, actor user.Actor) error
}

type DeleteUserService interface {
	Delete(ctx context.Context, idPengguna user.ID) error
}
