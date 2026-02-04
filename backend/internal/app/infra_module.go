package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outauth "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth_port_out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	outprofil "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type InfraModule struct {
	Pool *pgxpool.Pool
	Txm  txout.TxManager

	Sessions      out.SessionRepository
	AuthUsers     outauth.AuthUserrepository

	users         outuser.UserRepository
	profilSiswa   outuser.GetListSiswaRepo
	profilGuru    outuser.GetGuruListRepo
	
	profilSekolah outprofil.ProfilSekolahRepository
}

func BuildInfraModule(pool *pgxpool.Pool, logger corelog.Logger) *InfraModule {
	txm := pg.NewTxManager(pool, logger)

	return &InfraModule{
		Pool:          pool,
		Txm:           txm,
		Sessions:      pg.NewSessionRepo(pool, logger),
		AuthUsers:     pg.NewAuthUserRepo(pool, logger),
		users:         pg.NewUserRepo(pool, logger),
		profilSiswa:   pg.NewProfilSiswaRepo(pool, logger),
		profilGuru:    pg.NewProfilgGuruRepo(pool, logger),
		profilSekolah: pg.NewProfilSekolahRepo(pool, logger),
	}
}
