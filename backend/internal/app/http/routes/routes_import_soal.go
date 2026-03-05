package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterImportSoalRoutes(router *httprouter.Router, handlers ImportSoalHandlers, mw MiddlewareContract) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)

	router.POST("/admin/bank-soal/import/:idBankSoal", requireAdminGuru(mw.RateLimitStandard(handlers.ImportHandler.ImportSoal)))
	router.GET("/admin/bank-soal/import-job/:idJob", requireAdminGuru(mw.RateLimitStandard(handlers.GetJobHandler.GetJobByID)))
	router.GET("/admin/bank-soal/import-jobs/:idBankSoal", requireAdminGuru(mw.RateLimitStandard(handlers.GetJobHandler.GetJobsByBankSoal)))
}
