package user

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

type ID int
type Role string
type StatusAkun string
type Email string

const (
	ADMIN Role = "ADMIN"
	GURU  Role = "GURU"
	SISWA Role = "SISWA"
)

const (
	AKTIF    StatusAkun = "AKTIF"
	NONAKTIF StatusAkun = "NONAKTIF"
)

const (
	UsernameMinLength = 5
	UsernameMaxLength = 20
)

type Pengguna struct {
	ID             ID
	Username       string
	Email          *Email
	PasswordHashed string
	NamaLengkap    string
	JenisKelamin   string
	NoHp           *string
	Role           Role
	StatusAkun     StatusAkun
	Foto           string
}

type Actor struct {
	IdPengguna ID
	Username   string
	Role       Role
}

var (
	ErrInvalidEmail = errors.New("invalid Email")
)

func (role Role) ValidRole() bool {
	switch role {
	case ADMIN, GURU, SISWA:
		return true
	default:
		return false
	}
}

func CheckUsernameLength(username string) error {
	length := utf8.RuneCountInString(strings.TrimSpace(username))
	if length < UsernameMinLength || length > UsernameMaxLength {
		return coreerror.ErrUsernameLengthInvalid
	}
	return nil
}

func CheckNewEmail(raw *string) (Email, error) {
	if raw == nil {
		return "", ErrInvalidEmail
	}

	s := strings.TrimSpace(strings.ToLower(*raw))

	if s == "" {
		return "", ErrInvalidEmail
	}

	if len(s) > 254 {
		return "", ErrInvalidEmail
	}

	if _, err := mail.ParseAddress(s); err != nil {
		return "", ErrInvalidEmail
	}
	return Email(s), nil
}
