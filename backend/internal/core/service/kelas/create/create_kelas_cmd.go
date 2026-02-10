package kelas_service

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"

type CreateNamaKelasCmd struct {
	IdTingkatKelas kelas.ID
	NamaKelas      string
}

type CreateTingkatKelasCmd struct {
	TingkatKelas 	int
}