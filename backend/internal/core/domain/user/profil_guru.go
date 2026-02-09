package user

import (
	"errors"
	"strings"
	"unicode"
)

type NIP string



type ProfilGuru struct {
	ID          ID
	IdPengguna  ID
	Nip         NIP
	Jabatan     string
	BidangStudi string
}

var (
	ErrInvalidNIP = errors.New("invalid NIP")	
)

type DataGuru struct {
	IdPengguna 	 ID
	IdGuru	 	 ID
	Username     string
	Email        string
	NamaLengkap  string
	JenisKelamin string
	NoHp         string
	Foto         string
	Role         Role
	StatusAkun   StatusAkun

	Nip         string
	Jabatan     string
	BidangStudi string
	
}

func CheckNewNip(nip string) (NIP, error) {
	s := strings.TrimSpace(nip)

	if s == "-" {
		return NIP(s), nil
	}

	if len(s) != 18 {
		return "", ErrInvalidNIP
	}

	for _,r := range s {
		if !unicode.IsDigit(r) {
			return "", ErrInvalidNIP
		}
	}
	return NIP(s), nil
}
