package coreerror

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrDbError      = errors.New("db error")

	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrNoTokenProvided = errors.New("no token provided")
	ErrSessionExpired = errors.New("session expired")
	
	ErrNoFieldToUpdate = errors.New("no field to update")
	ErrUsernameTaken = errors.New("username taken")
	ErrNipTaken      = errors.New("NIP taken")
	ErrNisnTaken     = errors.New("NISN taken")
	ErrInvalidStatusAkun = errors.New("invalid status akun")

	ErrFileTooLarge = errors.New("file too large")
	ErrInvalidFileFormat = errors.New("invalid file format")
)
