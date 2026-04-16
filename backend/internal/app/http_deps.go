package app

import corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"

type HTTPDeps struct {
	Config Config
	Logger corelog.Logger

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
	Tokens        *TokenModule
	Infra         *InfraModule
}
