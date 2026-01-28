package user_service

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type CreateGuruCmd struct {
	Username     string
	Email        string
	Password     string
	NamaLengkap  string
	JenisKelamin string
	NoHp         string
	Foto         string

	Nip         string
	Jabatan     string
	BidangStudi string
}

type CreateSiswaCmd struct {
	Username     string
	Email        string
	Password     string
	NamaLengkap  string
	JenisKelamin string
	NoHp         string
	Foto         string

	IdTingkatKelas user.ID
	IdNamaKelas    user.ID
	Nisn           string
	NoAbsen        int
	Angkatan       int
	TempatLahir    string
	TanggalLahir   time.Time
}
