package user_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type ResetPasswordService struct {
	userRepo outuser.UserResetPasswordRepo
	hasher out.PasswordHasher
}


func NewResetPasswordService(userRepo outuser.UserResetPasswordRepo, hasher out.PasswordHasher) *ResetPasswordService {
	return &ResetPasswordService{userRepo: userRepo, hasher: hasher}
}

func(r *ResetPasswordService)ResetPasswordService(ctx context.Context, password string)error {
	logger := corelog.FromContext(ctx)

}
