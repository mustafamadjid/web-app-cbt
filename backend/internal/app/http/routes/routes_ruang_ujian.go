package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterRuangUjianRoutes(router *httprouter.Router, handlers RuangUjianHandlers, mw MiddlewareContract) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)

	router.GET("/admin/ruang-ujian", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetRuangUjian)))
	router.GET("/admin/ruang-ujian/id/:IdRuangan", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetRuangUjianByID)))
	router.GET("/admin/ruang-ujian/kode/:KodeRuang", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetRuangUjianByKode)))
	router.POST("/admin/ruang-ujian", requireAdminGuru(mw.RateLimitStandard(handlers.CreateHandler.CreateRuangUjian)))
	router.PATCH("/admin/ruang-ujian/:idRuangan", requireAdminGuru(mw.RateLimitStandard(handlers.UpdateHandler.UpdateRuangUjian)))
	router.DELETE("/admin/ruang-ujian/:idRuangan", requireAdminGuru(mw.RateLimitStandard(handlers.DeleteHandler.DeleteRuangUjian)))
}
