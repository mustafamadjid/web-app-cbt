package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterMataPelajaranRoutes(router *httprouter.Router, handlers MataPelajaranHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)

	router.GET("/admin/mata-pelajaran", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.ListMapel)))
	router.GET("/admin/mata-pelajaran/:idMapel", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetMapelByID)))
	router.POST("/admin/mata-pelajaran", requireAdmin(mw.RateLimitStandard(handlers.CreateHandler.CreateMapel)))
	router.PATCH("/admin/mata-pelajaran/:idMapel", requireAdmin(mw.RateLimitStandard(handlers.UpdateHandler.UpdateMapel)))
	router.DELETE("/admin/mata-pelajaran/:idMapel", requireAdmin(mw.RateLimitStandard(handlers.DeleteHandler.DeleteMapel)))
}
