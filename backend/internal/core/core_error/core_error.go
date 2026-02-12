package coreerror

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidInputSafe	= errors.New("invalid input : contains invalid characters")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrDbError      = errors.New("db error")
	ErrUserHasSession = errors.New("user has session")

	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrNoTokenProvided = errors.New("no token provided")
	ErrSessionExpired = errors.New("session expired")
	ErrHasSession = errors.New("user already has session")
	ErrNoSessionId = errors.New("no session id")
	ErrInvalidActionActivity = errors.New("invalid action activity")
	ErrInvalidIpAddress = errors.New("invalid ip address")

	ErrMissingId = errors.New("missing id")
	ErrMissingField = errors.New("missing field")
	
	ErrNoFieldToUpdate = errors.New("no field to update")
	ErrUsernameTaken = errors.New("username taken")
	ErrNipTaken      = errors.New("NIP taken")
	ErrNisnTaken     = errors.New("NISN taken")
	ErrInvalidStatusAkun = errors.New("invalid status akun")

	ErrTingkatKelasExist =errors.New("tingkat kelas already exist") 
	ErrNamaKelasExist =errors.New("nama kelas already exist") 
	ErrKodeMapelExist =errors.New("kode mapel already exist")

	ErrFileTooLarge = errors.New("file too large")
	ErrInvalidFileFormat = errors.New("invalid file format")

	ErrDeleteRestricted = errors.New("delete restricted")
)
