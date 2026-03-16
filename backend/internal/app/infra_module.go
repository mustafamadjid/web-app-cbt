package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	pgujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian"
	pgujianattempt "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/attempt"
	pgujianlist "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/list"
	pgujiansiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa"
	pgujiansiswaactiveattempt "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/active_attempt"
	pgujiansiswajawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/jawaban_ujian"
	pgujiansiswalist "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/list"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outaktivitas "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/aktivitas_user"
	outauth "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth_port_out"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	importsoal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	mapel_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/mata_pelajaran"
	pengumuman_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/pengumuman"
	outprofil "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
	ruangujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ruang_ujian"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
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

	kelasRepo              kelas_repo.KelasRepository
	bankSoalRepo           bank_soal_repo.BankSoalRepository
	mapelRepo              mapel_repo.MataPelajaranRepository
	pengumumanRepo         pengumuman_repo.PengumumanRepo
	ruangUjianRepo         ruangujian_repo.RuangUjianRepo
	listUjianRepo          ujian_repo.ListUjianRepository
	listUjianSiswaRepo     ujian_repo.ListUjianSiswaRepository
	soalUjianRepo          ujian_repo.SoalUjianRepository
	ujianRepo              ujian_repo.UjianRepository
	attemptUjianRepo       ujian_repo.AttemptUjianRepository
	activeAttemptUjianRepo ujian_repo.ActiveUjianAttemptRepository
	jawabanUjianRepo       ujian_repo.JawabanUjianRepository
	siswaUjianChecker      ujian_repo.SiswaUjianChecker
	waktuSelesaiUjianRepo  ujian_repo.WaktuSelesaiUjianRepository
	sesiRepo               sesi_repo.SesiRepository

	importSoalJobRepo importsoal_repo.ImportSoalJobRepo
	isiSoalBatchRepo  importsoal_repo.IsiSoalBatchRepo
}

func BuildInfraModule(pool *pgxpool.Pool, logger corelog.Logger) *InfraModule {
	txm := pg.NewTxManager(pool, logger)

	profilGuruRepo := pg.NewProfilgGuruRepo(pool, logger)
	profilSiswaRepo := pg.NewProfilSiswaRepo(pool, logger)
	ujianRepo := pgujian.NewUjianRepo(pool, logger, pool)
	siswaUjianRepo := pgujiansiswa.NewSiswaUjianRepo(pool, logger)
	listUjianSiswaRepo := pgujiansiswalist.NewListUjianSiswaRepo(pool, logger)
	attemptUjianRepo := pgujianattempt.NewAttemptUjianRepo(pool, logger)
	activeAttemptUjianRepo := pgujiansiswaactiveattempt.NewActiveAttemptRepo(pool, logger)
	jawabanUjianRepo := pgujiansiswajawaban.NewJawabanUjianRepo(pool, logger, pool)

	return &InfraModule{
		Pool:                   pool,
		Txm:                    txm,
		Sessions:               pg.NewSessionRepo(pool, logger),
		AuthUsers:              pg.NewAuthUserRepo(pool, logger),
		users:                  pg.NewUserRepo(pool, logger),
		userResetPasswordRepo:  pg.NewResetPasswordRepo(pool, logger),
		profilSiswa:            profilSiswaRepo,
		profilSiswaRepo:        profilSiswaRepo,
		profilGuru:             profilGuruRepo,
		profilGuruRepo:         profilGuruRepo,
		profilSekolah:          pg.NewProfilSekolahRepo(pool, logger),
		aktivitasUser:          pg.NewAktivitasUserRepo(pool, logger),
		kelasRepo:              pg.NewKelasRepo(pool, logger),
		bankSoalRepo:           pg.NewBankSoalRepo(pool, logger),
		mapelRepo:              pg.NewMapelRepo(pool, logger),
		pengumumanRepo:         pg.NewPengumumanRepo(pool, logger),
		ruangUjianRepo:         pg.NewRuangUjianRepo(pool, logger),
		listUjianRepo:          pgujianlist.NewListUjianRepo(pool, logger),
		listUjianSiswaRepo:     listUjianSiswaRepo,
		soalUjianRepo:          pgujianlist.NewListSoalUjianRepo(pool, logger),
		ujianRepo:              ujianRepo,
		attemptUjianRepo:       attemptUjianRepo,
		activeAttemptUjianRepo: activeAttemptUjianRepo,
		jawabanUjianRepo:       jawabanUjianRepo,
		siswaUjianChecker:      siswaUjianRepo,
		waktuSelesaiUjianRepo:  siswaUjianRepo,
		sesiRepo:               pg.NewSesirepo(pool, logger),
		importSoalJobRepo:      pg.NewImportSoalJobRepo(pool, logger),
		isiSoalBatchRepo:       pg.NewIsiSoalBatchRepo(pool, logger),
	}
}
