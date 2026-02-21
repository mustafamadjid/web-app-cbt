package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outaktivitas "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/aktivitas_user"
	outauth "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth_port_out"
	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	mapel_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/mata_pelajaran"
	pengumuman_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/pengumuman"
	outprofil "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
	ruangujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ruang_ujian"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type InfraModule struct {
	Pool *pgxpool.Pool
	Txm  txout.TxManager

	Sessions  out.SessionRepository
	AuthUsers outauth.AuthUserrepository

	users                 outuser.UserRepository
	userResetPasswordRepo outuser.UserResetPasswordRepo
	profilSiswa           outuser.GetListSiswaRepo
	profilSiswaRepo       outuser.ProfilSiswaRepository
	profilGuru            outuser.GetGuruListRepo
	profilGuruRepo        outuser.ProfilGuruRepository

	profilSekolah outprofil.ProfilSekolahRepository
	aktivitasUser outaktivitas.AktivitasUserRepository

	kelasRepo      kelas_repo.KelasRepository
	mapelRepo      mapel_repo.MataPelajaranRepository
	pengumumanRepo pengumuman_repo.PengumumanRepo
	ruangUjianRepo ruangujian_repo.RuangUjianRepo
	sesiRepo       sesi_repo.SesiRepository
}

func BuildInfraModule(pool *pgxpool.Pool, logger corelog.Logger) *InfraModule {
	txm := pg.NewTxManager(pool, logger)

	profilGuruRepo := pg.NewProfilgGuruRepo(pool, logger)
	profilSiswaRepo := pg.NewProfilSiswaRepo(pool, logger)

	return &InfraModule{
		Pool:                  pool,
		Txm:                   txm,
		Sessions:              pg.NewSessionRepo(pool, logger),
		AuthUsers:             pg.NewAuthUserRepo(pool, logger),
		users:                 pg.NewUserRepo(pool, logger),
		userResetPasswordRepo: pg.NewResetPasswordRepo(pool, logger),
		profilSiswa:           profilSiswaRepo,
		profilSiswaRepo:       profilSiswaRepo,
		profilGuru:            profilGuruRepo,
		profilGuruRepo:        profilGuruRepo,
		profilSekolah:         pg.NewProfilSekolahRepo(pool, logger),
		aktivitasUser:         pg.NewAktivitasUserRepo(pool, logger),
		kelasRepo:             pg.NewKelasRepo(pool, logger),
		mapelRepo:             pg.NewMapelRepo(pool, logger),
		pengumumanRepo:        pg.NewPengumumanRepo(pool, logger),
		ruangUjianRepo:        pg.NewRuangUjianRepo(pool, logger),
		sesiRepo:              pg.NewSesirepo(pool, logger),
	}
}
