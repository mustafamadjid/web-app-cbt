package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	pgaktivitasuser "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/aktivitas_user"
	pgauthuser "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/auth_user"
	pgbanksoal "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/bank_soal"
	pgimportsoaljob "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/import_soal_job"
	pgisisoalbatch "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/isi_soal_batch"
	pgkelas "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/kelas"
	pgmapel "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/mata_pelajaran"
	pgpengumuman "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/pengumuman"
	pgprofilguru "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_guru"
	pgprofilsekolah "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_sekolah"
	pgprofilsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_siswa"
	pgresetpassword "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/reset_password"
	pgruangujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ruang_ujian"
	pgsesi "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/sesi"
	pgsession "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/session"
	pgujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian"
	pgujianattempt "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/attempt"
	pgujiangrading "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/grading"
	pgujiangradinglist "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/grading/list"
	pgujianlist "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/list"
	pgujiansiswajawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/jawaban_ujian"
	pgujiansiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/ujian_siswa"
	pgujiansiswachecker "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/ujian_siswa_checker"
	pguser "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/user"
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
	grading_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/grading"
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

	kelasRepo         kelas_repo.KelasRepository
	bankSoalRepo      bank_soal_repo.BankSoalRepository
	mapelRepo         mapel_repo.MataPelajaranRepository
	pengumumanRepo    pengumuman_repo.PengumumanRepo
	ruangUjianRepo    ruangujian_repo.RuangUjianRepo
	listUjianRepo     ujian_repo.ListUjianRepository
	ujianSiswaRepo    ujian_repo.UjianSiswaRepository
	soalUjianRepo     ujian_repo.SoalUjianRepository
	ujianRepo         ujian_repo.UjianRepository
	attemptUjianRepo  ujian_repo.AttemptUjianRepository
	jawabanUjianRepo  ujian_repo.JawabanUjianRepository
	gradingRepo       grading_repo.GradingUjianRepository
	listGradingRepo   grading_repo.ListGradingUjianRepository
	gradingWorkerRepo grading_repo.GradingWorkerRepository
	siswaUjianChecker ujian_repo.SiswaUjianChecker
	sesiRepo          sesi_repo.SesiRepository

	importSoalJobRepo importsoal_repo.ImportSoalJobRepo
	isiSoalBatchRepo  importsoal_repo.IsiSoalBatchRepo
}

func BuildInfraModule(pool *pgxpool.Pool, logger corelog.Logger) *InfraModule {
	txm := pg.NewTxManager(pool, logger)

	profilGuruRepo := pgprofilguru.NewProfilgGuruRepo(pool, logger)
	profilSiswaRepo := pgprofilsiswa.NewProfilSiswaRepo(pool, logger)
	ujianRepo := pgujian.NewUjianRepo(pool, logger, pool)
	ujianSiswaRepo := pgujiansiswa.NewUjianSiswaRepo(pool, logger)
	siswaUjianChecker := pgujiansiswachecker.NewSiswaUjianCheckerRepo(pool, logger)
	attemptUjianRepo := pgujianattempt.NewAttemptUjianRepo(pool, logger)
	jawabanUjianRepo := pgujiansiswajawaban.NewJawabanUjianRepo(pool, logger)
	gradingRepo := pgujiangrading.NewGradingRepo(pool, logger)
	listGradingRepo := pgujiangradinglist.NewListGradingRepo(pool, logger)

	return &InfraModule{
		Pool:                  pool,
		Txm:                   txm,
		Sessions:              pgsession.NewSessionRepo(pool, logger),
		AuthUsers:             pgauthuser.NewAuthUserRepo(pool, logger),
		users:                 pguser.NewUserRepo(pool, logger),
		userResetPasswordRepo: pgresetpassword.NewResetPasswordRepo(pool, logger),
		profilSiswa:           profilSiswaRepo,
		profilSiswaRepo:       profilSiswaRepo,
		profilGuru:            profilGuruRepo,
		profilGuruRepo:        profilGuruRepo,
		profilSekolah:         pgprofilsekolah.NewProfilSekolahRepo(pool, logger),
		aktivitasUser:         pgaktivitasuser.NewAktivitasUserRepo(pool, logger),
		kelasRepo:             pgkelas.NewKelasRepo(pool, logger),
		bankSoalRepo:          pgbanksoal.NewBankSoalRepo(pool, logger),
		mapelRepo:             pgmapel.NewMapelRepo(pool, logger),
		pengumumanRepo:        pgpengumuman.NewPengumumanRepo(pool, logger),
		ruangUjianRepo:        pgruangujian.NewRuangUjianRepo(pool, logger),
		listUjianRepo:         pgujianlist.NewListUjianRepo(pool, logger),
		ujianSiswaRepo:        ujianSiswaRepo,
		soalUjianRepo:         pgujianlist.NewListSoalUjianRepo(pool, logger),
		ujianRepo:             ujianRepo,
		attemptUjianRepo:      attemptUjianRepo,
		jawabanUjianRepo:      jawabanUjianRepo,
		gradingRepo:           gradingRepo,
		listGradingRepo:       listGradingRepo,
		gradingWorkerRepo:     gradingRepo,
		siswaUjianChecker:     siswaUjianChecker,
		sesiRepo:              pgsesi.NewSesirepo(pool, logger),
		importSoalJobRepo:     pgimportsoaljob.NewImportSoalJobRepo(pool, logger),
		isiSoalBatchRepo:      pgisisoalbatch.NewIsiSoalBatchRepo(pool, logger),
	}
}
