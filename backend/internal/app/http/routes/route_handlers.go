package routes

import (
	authhttp "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/auth"
	aktivitasuserget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/aktivitas_user/get"
	banksoalcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/create"
	banksoaldelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/delete"
	banksoalget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/get"
	banksoalupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/update"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/import"
	kelascreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/create"
	kelasdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/delete"
	kelasget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/get"
	kelasupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/kelas/update"
	mapelcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/create"
	mapeldelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/delete"
	mapelget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/get"
	mapelupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/mata_pelajaran/update"
	pengumumancreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/create"
	pengumumandelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/delete"
	pengumumanget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/get"
	pengumumanupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/pengumuman/update"
	profilsekolahget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/get"
	profilsekolahupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah/update"
	ruangujiancreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/create"
	ruangujiandelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/delete"
	ruangujianget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/get"
	ruangujianupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ruang_ujian/update"
	sesicreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/create"
	sesidelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/delete"
	sesiget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/get"
	sesiupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/sesi/update"
	ujianget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/get"
	ujianlist "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list"
	usercreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/create"
	userdelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/delete"
	userget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/get"
	userresetpassword "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/reset_password"
	userupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/user/update"
)

type AuthHandlers struct {
	Handler *authhttp.AuthHandler
}

type UserHandlers struct {
	GetSiswaHandler *userget.GetSiswaHandler
	GetGuruHandler  *userget.GetGuruHandler
	CreateHandler   *usercreate.UserHandler
	UpdateHandler   *userupdate.UpdateHandler
	DeleteHandler   *userdelete.DeleteHandler
}

type ResetPasswordHandlers struct {
	Handler *userresetpassword.ResetPasswordHandler
}

type AktivitasUserHandlers struct {
	GetHandler *aktivitasuserget.AktivitasUserHandler
}

type ProfilSekolahHandlers struct {
	GetHandler    *profilsekolahget.GetProfilSekolahHandler
	UpdateHandler *profilsekolahupdate.UpdateProfilSekolahHandler
}

type KelasHandlers struct {
	GetHandler    *kelasget.GetKelasHandler
	CreateHandler *kelascreate.CreateKelasHandler
	UpdateHandler *kelasupdate.UpdateKelasHandler
	DeleteHandler *kelasdelete.DeleteKelasHandler
}

type MataPelajaranHandlers struct {
	GetHandler    *mapelget.GetMapelHandler
	CreateHandler *mapelcreate.CreateMapelHandler
	UpdateHandler *mapelupdate.UpdateMapelHandler
	DeleteHandler *mapeldelete.DeleteMapelHandler
}

type RuangUjianHandlers struct {
	GetHandler    *ruangujianget.GetRuangUjianHandler
	CreateHandler *ruangujiancreate.CreateRuangUjianHandler
	UpdateHandler *ruangujianupdate.UpdateRuangUjianHandler
	DeleteHandler *ruangujiandelete.DeleteRuangUjianHandler
}

type UjianHandlers struct {
	ListHandler *ujianlist.ListUjianHandler
	GetHandler  *ujianget.GetUjianHandler
}

type SesiHandlers struct {
	GetHandler    *sesiget.GetSesiHandler
	CreateHandler *sesicreate.CreateSesiHandler
	UpdateHandler *sesiupdate.UpdateSesiHandler
	DeleteHandler *sesidelete.DeleteSesiHandler
}

type PengumumanHandlers struct {
	GetHandler    *pengumumanget.GetPengumumanHandler
	CreateHandler *pengumumancreate.CreatePengumumanHandler
	UpdateHandler *pengumumanupdate.UpdatePengumumanHandler
	DeleteHandler *pengumumandelete.DeletePengumumanHandler
}

type ImportSoalHandlers struct {
	ImportHandler *importsoal.ImportHandler
	GetJobHandler *importsoal.GetJobHandler
}

type BankSoalHandlers struct {
	GetHandler    *banksoalget.GetBankSoalHandler
	CreateHandler *banksoalcreate.CreateBankSoalHandler
	UpdateHandler *banksoalupdate.UpdateBankSoalHandler
	DeleteHandler *banksoaldelete.DeleteBankSoalHandler
}
