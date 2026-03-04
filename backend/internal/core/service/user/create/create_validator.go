package user_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func validateCreateGuruActor(actor user.Actor) error {
	if actor.Role != user.ADMIN {
		return coreerror.ErrForbidden
	}
	return nil
}
func validateCreateSiswaActor(actor user.Actor) error {
	if actor.Role != user.ADMIN {
		return coreerror.ErrForbidden
	}
	return nil
}
