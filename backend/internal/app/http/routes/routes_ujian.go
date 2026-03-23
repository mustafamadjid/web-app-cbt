package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterUjianRoutes(router *httprouter.Router, handlers UjianHandlers, mw MiddlewareContract) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)
	requireAdmin := mw.RequireAccessRole(user.ADMIN)
	requireSiswa := mw.RequireAccessRole(user.SISWA)

	// Router siswa
	router.POST("/siswa/ujian/attempt", requireSiswa(mw.RateLimitStandard(handlers.AttemptUjianHandler.AttemptUjian)))
	router.GET("/siswa/ujian/attempt/active", requireSiswa(mw.RateLimitStandard(handlers.GetActiveAttemptUjianHandler.GetActiveAttemptUjian)))
	router.PATCH("/siswa/ujian/attempt/:idAttempt", requireSiswa(mw.RateLimitStandard(handlers.UpdateAttemptUjianHandler.UpdateAttemptUjian)))
	router.PATCH("/siswa/uijan/submit/:idAttempt", requireSiswa(mw.RateLimitStandard(handlers.SubmitUjianHandler.SubmitUjian)))
	router.GET("/siswa/ujian/jawaban/:idAttempt", requireSiswa(mw.RateLimitStandard(handlers.GetJawabanUjianHandler.GetJawabanUjian)))
	router.POST("/siswa/ujian/jawaban", requireSiswa(mw.RateLimitStandard(handlers.SaveJawabanUjianHandler.SaveJawabanUjian)))
	router.GET("/siswa/ujian/list", requireSiswa(mw.RateLimitStandard(handlers.ListUjianSiswaHandler.ListUjianSiswa)))
	router.GET("/siswa/ujian/waktu-selesai/:idJadwalUjian", requireSiswa(mw.RateLimitStandard(handlers.GetWaktuSelesaiUjianHandler.GetWaktuSelesaiUjian)))
	router.GET("/siswa/soal-ujian/:idJadwalUjian", requireSiswa(mw.RateLimitStandard(handlers.ListSoalUjianSiswaHandler.ListSoalUjianSiswa)))

	router.PATCH("/admin/ujian/attempt/:idAttempt/expire", requireAdmin(mw.RateLimitStandard(handlers.ExpireAttemptUjianHandler.ExpireAttemptUjian)))
	router.POST("/ujian", requireAdminGuru(mw.RateLimitStandard(handlers.CreateUjianHandler.CreateUjian)))
	router.GET("/jadwal-ujian", requireAdminGuru(mw.RateLimitStandard(handlers.ListHandler.ListUjian)))
	router.GET("/ujian/soal/bank-soal/:idBankSoal", requireAdminGuru(mw.RateLimitStandard(handlers.ListSoalUjianHandler.ListSoalUjian)))
	router.GET("/ujian/detail/:idUjian", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetUjianById)))

	router.PATCH("/ujian/detail/:idUjian", requireAdminGuru(mw.RateLimitStandard(handlers.UpdateUjianHandler.UpdateUjian)))
	router.DELETE("/ujian/detail/:idUjian", requireAdminGuru(mw.RateLimitStandard(handlers.DeleteUjianHandler.DeleteUjian)))
}
