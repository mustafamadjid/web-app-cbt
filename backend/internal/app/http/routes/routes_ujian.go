package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterUjianRoutes(router *httprouter.Router, handlers UjianHandlers, mw MiddlewareContract) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)

	router.GET("/jadwal-ujian", requireAdminGuru(mw.RateLimitStandard(handlers.ListHandler.ListUjian)))
	router.GET("/ujian/:idUjian", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetUjianByID)))
	router.GET("/jadwal-ujian/:idJadwalUjian", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetJadwalUjianByID)))
	router.GET("/peserta-ujian", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetAllPesertaUjian)))
	router.GET("/peserta-ujian/siswa/:idSiswa", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetPesertaUjianBySiswa)))
	router.GET("/jawaban-ujian-siswa", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetAllJawabanUjianSiswa)))
	router.GET("/jawaban-ujian-siswa/siswa/:idSiswa", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetJawabanBySiswa)))
}
