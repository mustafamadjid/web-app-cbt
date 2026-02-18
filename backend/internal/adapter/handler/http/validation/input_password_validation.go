package httpx

import "errors"

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("invalid password: password minimal 8 karakter")
	}

	return nil
}
