package out

import (
	"context"

	out_user "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type Tx interface {
	Pengguna() out_user.UserRepository
	ProfilGuru() out_user.ProfilGuruRepository
	ProfilSiswa() out_user.ProfilSiswaRepository

	Commit() error
	Rollback() error
}

type TxManager interface {
	Begin(ctx context.Context) (Tx, error)
}
