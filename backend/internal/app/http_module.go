package app

import (
	"net/http"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	httpmodule "github.com/mustafamadjid/web-app-cbt/internal/app/http"
	"github.com/mustafamadjid/web-app-cbt/internal/app/http/routes"
)

type HTTPModule struct {
	Handler http.Handler
	Server  *http.Server
}

func BuildHTTPModule(deps HTTPDeps) *HTTPModule {
	cookies := cookie.CookieConfig{
		AccessName:  deps.Config.Cookie.AccessName,
		RefreshName: deps.Config.Cookie.RefreshName,
		Domain:      deps.Config.Cookie.Domain,
		Secure:      deps.Config.Cookie.Secure,
		SameSite:    deps.Config.Cookie.SameSite,
	}

	module := httpmodule.BuildHTTPModule(httpmodule.HTTPDeps{
		HTTPAddr:           deps.Config.HTTP.Addr,
		CookieConfig:       cookies,
		ImageStoreDir:      deps.Config.ImageStore.Dir,
		ImageStoreRoute:    deps.Config.ImageStore.Route,
		DocumentStoreDir:   deps.Config.DocumentStore.Dir,
		DocumentStoreRoute: deps.Config.DocumentStore.Route,
		Logger:             deps.Logger,
		AccessTokenSvc:     deps.Tokens.AccessTokenSvc,
		RefreshTokenSvc:    deps.Tokens.RefreshTokenSvc,
		Sessions:           deps.Infra.Sessions,
		Auth: routes.AuthHandlers{
			Handler: deps.Auth.Handler,
		},
		Users: routes.UserHandlers{
			GetSiswaHandler: deps.Users.GetSiswaHandler,
			GetGuruHandler:  deps.Users.GetGuruHandler,
			CreateHandler:   deps.Users.CreateHandler,
			UpdateHandler:   deps.Users.UpdateHandler,
			DeleteHandler:   deps.Users.DeleteHandler,
		},
		ResetPassword: routes.ResetPasswordHandlers{
			Handler: deps.ResetPassword.Handler,
		},
		AktivitasUser: routes.AktivitasUserHandlers{
			GetHandler: deps.AktivitasUser.GetHandler,
		},
		ProfilSekolah: routes.ProfilSekolahHandlers{
			GetHandler:    deps.ProfilSekolah.GetHandler,
			UpdateHandler: deps.ProfilSekolah.UpdateHandler,
		},
		Kelas: routes.KelasHandlers{
			GetHandler:    deps.Kelas.GetHandler,
			CreateHandler: deps.Kelas.CreateHandler,
			UpdateHandler: deps.Kelas.UpdateHandler,
			DeleteHandler: deps.Kelas.DeleteHandler,
		},
		MataPelajaran: routes.MataPelajaranHandlers{
			GetHandler:    deps.MataPelajaran.GetHandler,
			CreateHandler: deps.MataPelajaran.CreateHandler,
			UpdateHandler: deps.MataPelajaran.UpdateHandler,
			DeleteHandler: deps.MataPelajaran.DeleteHandler,
		},
		RuangUjian: routes.RuangUjianHandlers{
			GetHandler:    deps.RuangUjian.GetHandler,
			CreateHandler: deps.RuangUjian.CreateHandler,
			UpdateHandler: deps.RuangUjian.UpdateHandler,
			DeleteHandler: deps.RuangUjian.DeleteHandler,
		},
		Ujian: routes.UjianHandlers{
			AttemptUjianHandler:   deps.Ujian.AttemptUjianHandler,
			CreateUjianHandler:    deps.Ujian.CreateUjianHandler,
			ListHandler:           deps.Ujian.ListHandler,
			ListUjianSiswaHandler: deps.Ujian.ListUjianSiswaHandler,
			ListSoalUjianHandler:  deps.Ujian.ListSoalUjianHandler,
			GetHandler:            deps.Ujian.GetHandler,
			UpdateUjianHandler:    deps.Ujian.UpdateUjianHandler,
			DeleteUjianHandler:    deps.Ujian.DeleteUjianHandler,
		},
		Sesi: routes.SesiHandlers{
			GetHandler:    deps.Sesi.GetHandler,
			CreateHandler: deps.Sesi.CreateHandler,
			UpdateHandler: deps.Sesi.UpdateHandler,
			DeleteHandler: deps.Sesi.DeleteHandler,
		},
		Pengumuman: routes.PengumumanHandlers{
			GetHandler:    deps.Pengumuman.GetHandler,
			CreateHandler: deps.Pengumuman.CreateHandler,
			UpdateHandler: deps.Pengumuman.UpdateHandler,
			DeleteHandler: deps.Pengumuman.DeleteHandler,
		},
		ImportSoal: routes.ImportSoalHandlers{
			ImportHandler: deps.ImportSoal.ImportHandler,
			GetJobHandler: deps.ImportSoal.GetJobHandler,
		},
		BankSoal: routes.BankSoalHandlers{
			GetHandler:    deps.BankSoal.GetHandler,
			CreateHandler: deps.BankSoal.CreateHandler,
			UpdateHandler: deps.BankSoal.UpdateHandler,
			DeleteHandler: deps.BankSoal.DeleteHandler,
		},
	})

	return &HTTPModule{
		Handler: module.Handler,
		Server:  module.Server,
	}
}
