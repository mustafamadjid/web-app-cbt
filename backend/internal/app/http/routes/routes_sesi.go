package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterSesiRoutes(router *httprouter.Router, handlers SesiHandlers, mw MiddlewareContract) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)

	router.GET("/admin/sesi/login-aktif", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.ListActiveLoginSession)))
	router.GET("/admin/sesi", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.ListSesi)))
	router.GET("/admin/sesi/sesi-id/:idSesi", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetSesiByID)))
	router.GET("/admin/sesi/kode/:kodeSesi", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetSesiByKode)))
	router.POST("/admin/sesi", requireAdminGuru(mw.RateLimitStandard(handlers.CreateHandler.CreateSesiHandler)))
	router.PATCH("/admin/sesi/:idSesi", requireAdminGuru(mw.RateLimitStandard(handlers.UpdateHandler.UpdateSesi)))
	router.DELETE("/admin/sesi/:idSesi", requireAdminGuru(mw.RateLimitStandard(handlers.DeleteHandler.DeleteSesi)))
}
