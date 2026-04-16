package routes

import (
	authhttp "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/auth"
	aktivitasuserget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/aktivitas_user/get"
	banksoalcreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/create"
	banksoaldelete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/delete"
	banksoalget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/get"
	banksoalupdate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/update"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/import"
	dashboardget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/dashboard/get"
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
	ujianactiveattempt "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/active_attempt"
	ujiancreateattempt "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/create"
	ujianexpireattempt "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/expire"
	ujianupdateattempt "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/update"
	ujianlistpesertasubmitted "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/list_peserta_submitted"
	ujiangetjawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/jawaban_ujian/get_jawaban"
	ujianhasiljawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/jawaban_ujian/hasil_jawaban"
	ujiansavejawaban "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/jawaban_ujian/save_jawaban"
	ujiankoreksiessay "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/koreksi_essay"
	ujianlist "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list"
	ujianlistsoal "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian"
	ujianlistsoalsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_soal_ujian_for_siswa"
	ujianlistessayungraded "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_ujian_essay_ungraded"
	ujianlistselesai "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_ujian_selesai_siswa"
	ujianlistsiswa "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/list_ujian_siswa"
	ujianstatistik "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/statistik_ujian"
	ujiansubmit "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/submit_ujian"
	ujiancreate "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/create"
	ujiandeleteujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/delete"
	ujianget "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/get"
	ujianupdateujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/ujian_siswa/update"
	ujiangetwaktuujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/waktu_ujian"
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

type DashboardHandlers struct {
	GetHandler *dashboardget.GetDashboardStatistikHandler
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
	AttemptUjianHandler              *ujiancreateattempt.AttemptUjianHandler
	GetActiveAttemptUjianHandler     *ujianactiveattempt.GetActiveAttemptUjianHandler
	ListPesertaUjianSubmittedHandler *ujianlistpesertasubmitted.ListPesertaUjianSubmittedHandler
	UpdateAttemptUjianHandler        *ujianupdateattempt.UpdateAttemptUjianHandler
	ExpireAttemptUjianHandler        *ujianexpireattempt.ExpireAttemptUjianHandler
	SubmitUjianHandler               *ujiansubmit.SubmitUjianHandler
	GetJawabanUjianHandler           *ujiangetjawaban.GetJawabanUjianHandler
	HasilJawabanUjianHandler         *ujianhasiljawaban.HasilJawabanUjianHandler
	SaveJawabanUjianHandler          *ujiansavejawaban.SaveJawabanUjianHandler
	KoreksiEssayHandler              *ujiankoreksiessay.KoreksiEssayHandler
	ListUjianEssayUngradedHandler    *ujianlistessayungraded.ListUjianEssayUngradedHandler
	ListUjianSelesaiSiswaHandler     *ujianlistselesai.ListUjianSelesaiSiswaHandler
	GetStatistikUjianHandler         *ujianstatistik.GetStatistikUjianHandler
	CreateUjianHandler               *ujiancreate.CreateRuangUjianHandler
	ListHandler                      *ujianlist.ListUjianHandler
	ListUjianSiswaHandler            *ujianlistsiswa.ListUjianSiswaHandler
	ListSoalUjianHandler             *ujianlistsoal.ListSoalUjianHandler
	ListSoalUjianSiswaHandler        *ujianlistsoalsiswa.ListSoalUjianSiswaHandler
	GetWaktuSelesaiUjianHandler      *ujiangetwaktuujian.GetWaktuSelesaiUjianHandler
	GetHandler                       *ujianget.GetUjianHandler
	UpdateUjianHandler               *ujianupdateujian.UpdateUjianHandler
	DeleteUjianHandler               *ujiandeleteujian.DeleteUjianHandler
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
