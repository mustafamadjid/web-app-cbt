package user_service

import (
	"errors"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func validateUpdateUserActor(actor user.Actor) error {
	if actor.Role != user.ADMIN {
		return coreerror.ErrForbidden
	}
	return nil
}
func validateUpdateUserID(idPengguna user.ID) error {
	if idPengguna == 0 {
		return errors.New("Id pengguna required")
	}
	return nil
}
