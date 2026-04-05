package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterSesiRoutes(router *httprouter.Router, handlers SesiHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)

	router.GET("/admin/sesi/login-aktif", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.ListActiveLoginSession)))
	router.GET("/admin/sesi", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.ListSesi)))
	router.GET("/admin/sesi/sesi-id/:idSesi", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetSesiByID)))
	router.GET("/admin/sesi/kode/:kodeSesi", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetSesiByKode)))
	router.POST("/admin/sesi", requireAdmin(mw.RateLimitStandard(handlers.CreateHandler.CreateSesiHandler)))
	router.PATCH("/admin/sesi/:idSesi", requireAdmin(mw.RateLimitStandard(handlers.UpdateHandler.UpdateSesi)))
	router.DELETE("/admin/sesi/:idSesi", requireAdmin(mw.RateLimitStandard(handlers.DeleteHandler.DeleteSesi)))
}
