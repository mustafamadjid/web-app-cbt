package out

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type ProfilSiswaRepository interface {
	FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error)
	ExistByNISN(ctx context.Context, nisn string) (bool, error)
	CreateProfilSiswa(ctx context.Context, profilSiswa user.ProfilSiswa) (user.ID, error)

	UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa updatepatch.ProfilSiswa) error
}

type GetListSiswaRepo interface {
	GetListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem, error)
}
