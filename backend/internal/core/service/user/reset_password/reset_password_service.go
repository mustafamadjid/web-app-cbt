package user_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type ResetPasswordService struct {
	userRepo outuser.UserResetPasswordRepo
	hasher   out.PasswordHasher
}

func NewResetPasswordService(userRepo outuser.UserResetPasswordRepo, hasher out.PasswordHasher) *ResetPasswordService {
	return &ResetPasswordService{userRepo: userRepo, hasher: hasher}
}

func (r *ResetPasswordService) ResetPasswordService(ctx context.Context, userID user.ID, password string) error {
	logger := corelog.FromContext(ctx)

	hashedPassword, err := r.hasher.GenerateHash(password)
	if err != nil {
		logger.Error(ctx, "failed hashing password", "layer", "core.service", "op", "user.reset_password", "err", err)
		return err
	}

	if err := r.userRepo.ResetPassword(ctx, userID, hashedPassword); err != nil {
		logger.Error(ctx, "failed reset password", "layer", "core.service", "op", "user.reset_password", "err", err)
		return err
	}
	return nil

}
