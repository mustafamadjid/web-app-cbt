package query

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type SiswaListItem struct {
	IdPengguna   user.ID
	Username     string
	Email        user.Email
	NamaLengkap  string
	JenisKelamin string
	NoHp         string
	Foto         string
	NoAbsen      int
	StatusAkun   user.StatusAkun
	NamaKelas    string
	TingkatKelas int
	Angkatan     int
	TempatLahir  string
	TanggalLahir time.Time
}

type ListSiswaFilter struct {
	Search   string
	Status   *user.StatusAkun
	Limit    int
	Offset   int
	SortBy   string
	SortDesc bool

	Angkatan     *int
	TingkatKelas *int
	JenisKelamin *int
}
