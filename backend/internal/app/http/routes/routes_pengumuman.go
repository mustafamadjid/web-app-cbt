package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterPengumumanRoutes(router *httprouter.Router, handlers PengumumanHandlers, mw MiddlewareContract) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)

	router.GET("/pengumuman/active", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetPengumumanActive)))
	router.GET("/pengumuman/non-active", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetPengumumanNonActive)))
	router.GET("/pengumuman/incoming", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetPengumumanIncoming)))
	router.GET("/pengumuman/id/:idPengumuman", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetPengumumanByID)))
	router.POST("/pengumuman", requireAdminGuru(mw.RateLimitStandard(handlers.CreateHandler.CreatePengumuman)))
	router.PATCH("/pengumuman/:idPengumuman", requireAdminGuru(mw.RateLimitStandard(handlers.UpdateHandler.UpdatePengumuman)))
	router.DELETE("/pengumuman/:idPengumuman", requireAdminGuru(mw.RateLimitStandard(handlers.DeleteHandler.DeletePengumuman)))
}
