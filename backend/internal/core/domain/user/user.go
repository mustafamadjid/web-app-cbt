package user

import (
	"errors"
	"net/mail"
	"strings"
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

type Pengguna struct {
	ID             ID
	Username       string
	Email          Email
	PasswordHashed string
	NamaLengkap    string
	JenisKelamin   string
	NoHp           string
	Role           Role
	StatusAkun     StatusAkun
	Foto           string
}

type Actor struct {
	IdPengguna ID
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

func CheckNewEmail(raw string) (Email, error) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if s == "" {
		return "", ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(s); err != nil {
		return "", ErrInvalidEmail
	}
	return Email(s), nil
}