package user_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type ValidatedGuruCreate struct {
	cmd          CreateGuruCmd
	email        *user.Email
	nipToStore   user.NIP
	nipValidated user.NIP
	isDashedNip  bool
}

type ValidatedSiswaCreate struct {
	cmd           CreateSiswaCmd
	email         *user.Email
	nisnToStore   user.NISN
	nisnValidated user.NISN
	isDashedNisn  bool
}

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

func validateCreateEmail(raw *string) (*user.Email, error) {
	if raw == nil {
		return nil, nil
	}

	email, err := user.CheckNewEmail(raw)
	if err != nil {
		return nil, err
	}
	return &email, nil
}
