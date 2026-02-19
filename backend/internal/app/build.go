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
	Kelas         *KelasModule
	MataPelajaran *MataPelajaranModule
	RuangUjian    *RuangUjianModule
	Sesi          *SesiModule
	ResetPassword *ResetPasswordModule
	HTTP          *HTTPModule
}

func Build(ctx context.Context, cfg Config, dbURL string, hasher out.PasswordHasher, logger corelog.Logger) (*App, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	infra := BuildInfraModule(pool, logger)
	tokens := BuildTokenModule(cfg)
	aktivitasUser := BuildAktivitasUserModule(infra)
	auth := BuildAuthModule(cfg, infra, tokens, hasher, aktivitasUser)
	users := BuildUserModule(cfg, infra, hasher, aktivitasUser)
	profilSekolah := BuildProfilSekolahModule(cfg, infra)
	kelas := BuildKelasModule(infra, aktivitasUser)
	mapel := BuildMataPelajaranModule(infra)
	ruangUjian := BuildRuangUjianModule(infra)
	sesi := BuildSesiModule(infra)
	resetPassword := BuildResetPasswordModule(infra, hasher)
	httpm := BuildHTTPModule(cfg, auth, users, profilSekolah, aktivitasUser, kelas, mapel, ruangUjian, sesi, resetPassword, tokens, infra, logger)

	return &App{
		Infra:         infra,
		Tokens:        tokens,
		Auth:          auth,
		Users:         users,
		ProfilSekolah: profilSekolah,
		AktivitasUser: aktivitasUser,
		Kelas:         kelas,
		MataPelajaran: mapel,
		RuangUjian:    ruangUjian,
		Sesi:          sesi,
		ResetPassword: resetPassword,
		HTTP:          httpm,
	}, nil
}
