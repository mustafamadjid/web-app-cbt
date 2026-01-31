package out

import (
	"context"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
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

	UpdateProfilSiswa(ctx context.Context,idPengguna user.ID,profilSiswa UpdateProfilSiswaPatch) error
}

type GetListSiswaRepo interface {
	GetListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem,error)
}