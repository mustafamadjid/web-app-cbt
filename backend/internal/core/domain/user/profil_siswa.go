package user

import (
	"errors"
	"strings"
	"time"
	"unicode"
)

type NISN string

type ProfilSiswa struct {
	ID             ID
	IdPengguna     ID
	IdTingkatKelas ID
	IdNamaKelas    ID
	Nisn           NISN
	NoAbsen        int
	Angkatan       int
	TempatLahir    string
	TanggalLahir   time.Time
}

var (
	ErrInvalidNISN = errors.New("invalid NISN")
	ErrInvalidAbsen = errors.New("invalid absen")
	ErrInvalidAngkatan = errors.New("invalid Angkatan")
)

func CheckNewNISN(nisn string) (NISN, error) {
	s := strings.TrimSpace(nisn)
	if len(s) != 10 {
		return "", ErrInvalidNISN
	}

	for _, r := range s {
		if !unicode.IsDigit(r) {
			return "", ErrInvalidNISN
		}
	}
	return NISN(s), nil
}

func CheckAbsen(noAbsen int) error {
	if noAbsen < 1 {
		return ErrInvalidAbsen
	}
	return nil
}

func CheckAngkatan(angkatan int) error {
	if angkatan < 2019 || angkatan > time.Now().Year() {
		return ErrInvalidAngkatan
	} 
	return nil
}