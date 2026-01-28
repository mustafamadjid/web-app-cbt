package out

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type UpdateProfilGuruPatch struct {
	Nip        *string
	Jabatan    *string
	BidangStudi *string
}


type ProfilGuruRepository interface {
	FindProfilGuruByID(ctx context.Context, id user.ID) (user.ProfilGuru, error)
	ExistByNIP(ctx context.Context, nip user.NIP) (bool, error)
	CreateProfilGuru(ctx context.Context, profilGuru user.ProfilGuru) (user.ID, error)
	UpdateProfilGuru(ctx context.Context,idPengguna user.ID, profilGuru UpdateProfilGuruPatch) error
}