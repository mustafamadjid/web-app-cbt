package httpmodule

import (
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/app/http/routes"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type HTTPDeps struct {
	HTTPAddr string

	CookieConfig cookie.CookieConfig

	ImageStoreDir      string
	ImageStoreRoute    string
	DocumentStoreDir   string
	DocumentStoreRoute string

	Logger          corelog.Logger
	AccessTokenSvc  out.AccessTokenService
	RefreshTokenSvc out.RefreshTokenService
	Sessions        out.SessionRepository

	Auth          routes.AuthHandlers
	Users         routes.UserHandlers
	ResetPassword routes.ResetPasswordHandlers
	AktivitasUser routes.AktivitasUserHandlers
	ProfilSekolah routes.ProfilSekolahHandlers
	Kelas         routes.KelasHandlers
	MataPelajaran routes.MataPelajaranHandlers
	RuangUjian    routes.RuangUjianHandlers
	Ujian         routes.UjianHandlers
	Sesi          routes.SesiHandlers
	Pengumuman    routes.PengumumanHandlers
	ImportSoal    routes.ImportSoalHandlers
	BankSoal      routes.BankSoalHandlers
}
