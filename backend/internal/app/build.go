package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type App struct {
	Infra         *InfraModule
	Tokens        *TokenModule
	Auth          *AuthModule
	Users         *UserModule
	ProfilSekolah *ProfilSekolahModule
	AktivitasUser *AktivitasUserModule
	Dashboard     *DashboardModule
	Kelas         *KelasModule
	MataPelajaran *MataPelajaranModule
	RuangUjian    *RuangUjianModule
	Ujian         *UjianModule
	Sesi          *SesiModule
	Pengumuman    *PengumumanModule
	BankSoal      *BankSoalModule
	ResetPassword *ResetPasswordModule
	ImportSoal    *ImportSoalModule
	HTTP          *HTTPModule
}

func Build(ctx context.Context, cfg Config, dbURL string, hasher out.PasswordHasher, logger corelog.Logger, deleteFile *DeleteFileModule) (*App, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	infra := BuildInfraModule(pool, logger)
	tokens := BuildTokenModule(cfg)
	aktivitasUser := BuildAktivitasUserModule(infra)
	auth := BuildAuthModule(cfg, infra, tokens, hasher, aktivitasUser)
	users := BuildUserModule(cfg, infra, hasher, aktivitasUser, deleteFile)
	profilSekolah := BuildProfilSekolahModule(cfg, infra)
	dashboard := BuildDashboardModule(infra)
	kelas := BuildKelasModule(infra, aktivitasUser)
	mapel := BuildMataPelajaranModule(infra)
	ruangUjian := BuildRuangUjianModule(infra)
	ujian := BuildUjianModule(infra)
	sesi := BuildSesiModule(infra)
	pengumuman := BuildPengumumanModule(cfg, infra, deleteFile)
	bankSoal := BuildBankSoalModule(infra)
	resetPassword := BuildResetPasswordModule(infra, hasher)
	importSoal := BuildImportSoalModule(infra, cfg, logger)
	httpm := BuildHTTPModule(HTTPDeps{
		Config:        cfg,
		Logger:        logger,
		Auth:          auth,
		Users:         users,
		ProfilSekolah: profilSekolah,
		AktivitasUser: aktivitasUser,
		Dashboard:     dashboard,
		Kelas:         kelas,
		MataPelajaran: mapel,
		RuangUjian:    ruangUjian,
		Ujian:         ujian,
		Sesi:          sesi,
		Pengumuman:    pengumuman,
		BankSoal:      bankSoal,
		ResetPassword: resetPassword,
		ImportSoal:    importSoal,
		Tokens:        tokens,
		Infra:         infra,
	})

	return &App{
		Infra:         infra,
		Tokens:        tokens,
		Auth:          auth,
		Users:         users,
		ProfilSekolah: profilSekolah,
		AktivitasUser: aktivitasUser,
		Dashboard:     dashboard,
		Kelas:         kelas,
		MataPelajaran: mapel,
		RuangUjian:    ruangUjian,
		Ujian:         ujian,
		Sesi:          sesi,
		Pengumuman:    pengumuman,
		BankSoal:      bankSoal,
		ResetPassword: resetPassword,
		ImportSoal:    importSoal,
		HTTP:          httpm,
	}, nil
}
