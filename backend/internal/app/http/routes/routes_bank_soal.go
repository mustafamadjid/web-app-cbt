package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterBankSoalRoutes(router *httprouter.Router, handlers BankSoalHandlers, mw MiddlewareContract) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)

	router.GET("/admin-guru/bank-soal", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetBankSoal)))
	router.GET("/admin-guru/bank-soal-uploaded", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetBankSoalUploaded)))
	router.POST("/admin-guru/bank-soal", requireAdminGuru(mw.RateLimitStandard(handlers.CreateHandler.CreateBankSoal)))
	router.GET("/admin-guru/guru/bank-soal/:idPengguna", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetBankSoalByGuru)))
	router.GET("/admin-guru/bank-soal/:idBankSoal", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetBankSoalByID)))
	router.PATCH("/admin-guru/bank-soal/:idBankSoal", requireAdminGuru(mw.RateLimitStandard(handlers.UpdateHandler.UpdateBankSoal)))
	router.DELETE("/admin-guru/bank-soal/:idBankSoal", requireAdminGuru(mw.RateLimitStandard(handlers.DeleteHandler.DeleteBankSoal)))
}
