package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterPengumumanRoutes(router *httprouter.Router, handlers PengumumanHandlers, mw MiddlewareContract) {
	requireAdminGuruSiswa := mw.RequireAccessRole(user.ADMIN, user.GURU,user.SISWA)

	router.GET("/pengumuman/active", requireAdminGuruSiswa(mw.RateLimitStandard(handlers.GetHandler.GetPengumumanActive)))
	router.GET("/pengumuman/non-active", requireAdminGuruSiswa(mw.RateLimitStandard(handlers.GetHandler.GetPengumumanNonActive)))
	router.GET("/pengumuman/incoming", requireAdminGuruSiswa(mw.RateLimitStandard(handlers.GetHandler.GetPengumumanIncoming)))
	router.GET("/pengumuman/id/:idPengumuman", requireAdminGuruSiswa(mw.RateLimitStandard(handlers.GetHandler.GetPengumumanByID)))
	router.POST("/pengumuman", requireAdminGuruSiswa(mw.RateLimitStandard(handlers.CreateHandler.CreatePengumuman)))
	router.PATCH("/pengumuman/:idPengumuman", requireAdminGuruSiswa(mw.RateLimitStandard(handlers.UpdateHandler.UpdatePengumuman)))
	router.DELETE("/pengumuman/:idPengumuman", requireAdminGuruSiswa(mw.RateLimitStandard(handlers.DeleteHandler.DeletePengumuman)))
}
