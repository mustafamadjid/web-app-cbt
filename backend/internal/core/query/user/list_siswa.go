package query

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type SiswaListItem struct {
	IdPengguna user.ID
	Username  		string
	Email      		user.Email
	NamaLengkap 	string
	JenisKelamin 	string
	NoHp        	string
	Foto 			string
	StatusAkun  	user.StatusAkun
	NamaKelas 		string
	TingkatKelas 	int
	Angkatan 		int
}

type ListSiswaFilter struct {
	Search    string
	Status    *user.StatusAkun
	Limit     int
	Offset    int
	SortBy    string
	SortDesc  bool

	Angkatan *int
	TingkatKelas *int
	JenisKelamin *string
}