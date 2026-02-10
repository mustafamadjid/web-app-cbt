package out

import (
	"context"

	
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type Tx interface {
	Pengguna() out.UserRepository
	ProfilGuru() out.ProfilGuruRepository
	ProfilSiswa() out.ProfilSiswaRepository

	Commit() error
	Rollback() error
}

type TxManager interface {
	Begin(ctx context.Context) (Tx, error)
}
