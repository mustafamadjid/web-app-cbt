package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outaktivitas "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/aktivitas_user"
	outauth "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth_port_out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	outprofil "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type InfraModule struct {
	Pool *pgxpool.Pool
	Txm  txout.TxManager

	Sessions  out.SessionRepository
	AuthUsers outauth.AuthUserrepository

	users          outuser.UserRepository
	profilSiswa    outuser.GetListSiswaRepo
	profilGuru     outuser.GetGuruListRepo
	profilGuruRepo outuser.ProfilGuruRepository

	profilSekolah outprofil.ProfilSekolahRepository
	aktivitasUser outaktivitas.AktivitasUserRepository
}

func BuildInfraModule(pool *pgxpool.Pool, logger corelog.Logger) *InfraModule {
	txm := pg.NewTxManager(pool, logger)

	profilGuruRepo := pg.NewProfilgGuruRepo(pool, logger)

	return &InfraModule{
		Pool:           pool,
		Txm:            txm,
		Sessions:       pg.NewSessionRepo(pool, logger),
		AuthUsers:      pg.NewAuthUserRepo(pool, logger),
		users:          pg.NewUserRepo(pool, logger),
		profilSiswa:    pg.NewProfilSiswaRepo(pool, logger),
		profilGuru:     profilGuruRepo,
		profilGuruRepo: profilGuruRepo,
		profilSekolah:  pg.NewProfilSekolahRepo(pool, logger),
		aktivitasUser:  pg.NewAktivitasUserRepo(pool, logger),
	}
}
