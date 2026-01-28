package coreerror

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidCreds = errors.New("invalid credentials")

	ErrUsernameTaken = errors.New("username taken")
	ErrNipTaken      = errors.New("NIP taken")
	ErrNisnTaken     = errors.New("NISN taken")
)
