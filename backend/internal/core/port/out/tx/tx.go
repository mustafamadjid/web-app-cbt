package out

import (
	"context"

	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type Tx interface {
	Pengguna() out.UserRepository
	ProfilGuru() out.ProfilGuruRepository
	ProfilSiswa() out.ProfilSiswaRepository

	Kelas() kelas_repo.KelasRepository

	Commit() error
	Rollback() error
}

type TxManager interface {
	Begin(ctx context.Context) (Tx, error)
}
