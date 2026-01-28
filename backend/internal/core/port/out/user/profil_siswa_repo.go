package out

import (
	"context"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type UpdateProfilSiswaPatch struct {
	IdTingkatKelas *user.ID
	IdNamaKelas    *user.ID
	Nisn           *string
	NoAbsen        *int
	Angkatan       *int
	TempatLahir    *string
	TanggalLahir   *time.Time
}


type ProfilSiswaRepository interface {
	FindProfilSiswaByID(ctx context.Context,id user.ID) (user.ProfilSiswa,error)
	ExistByNISN(ctx context.Context,nisn string) (bool,error)
	CreateProfilSiswa(ctx context.Context,profilSiswa user.ProfilSiswa) (user.ID,error)
	UpdateProfilSiswa(ctx context.Context,profilSiswa user.ProfilSiswa) (user.ID,error)
}