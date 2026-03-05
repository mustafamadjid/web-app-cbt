package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterRuangUjianRoutes(router *httprouter.Router, handlers RuangUjianHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)

	router.GET("/admin/ruang-ujian", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetRuangUjian)))
	router.GET("/admin/ruang-ujian/id/:IdRuangan", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetRuangUjianByID)))
	router.GET("/admin/ruang-ujian/kode/:KodeRuang", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetRuangUjianByKode)))
	router.POST("/admin/ruang-ujian", requireAdmin(mw.RateLimitStandard(handlers.CreateHandler.CreateRuangUian)))
	router.PATCH("/admin/ruang-ujian/:idRuangan", requireAdmin(mw.RateLimitStandard(handlers.UpdateHandler.UpdateRuangUjian)))
	router.DELETE("/admin/ruang-ujian/:idRuangan", requireAdmin(mw.RateLimitStandard(handlers.DeleteHandler.DeleteRuangUjian)))
}
