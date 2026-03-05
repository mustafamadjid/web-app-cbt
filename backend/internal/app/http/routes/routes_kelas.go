package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterKelasRoutes(router *httprouter.Router, handlers KelasHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)

	router.GET("/admin/kelas", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.ListKelas)))
	router.GET("/admin/kelas/:idTingkatKelas/:idNamaKelas", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetKelasByID)))
	router.POST("/admin/kelas/tingkat-kelas", requireAdmin(mw.RateLimitStandard(handlers.CreateHandler.CreateTingkatKelas)))
	router.POST("/admin/kelas/nama-kelas", requireAdmin(mw.RateLimitStandard(handlers.CreateHandler.CreateNamaKelas)))
	router.PATCH("/admin/kelas/nama-kelas/:idNamaKelas", requireAdmin(mw.RateLimitStandard(handlers.UpdateHandler.UpdateNamaKelas)))
	router.DELETE("/admin/kelas/nama-kelas/:idNamaKelas", requireAdmin(mw.RateLimitStandard(handlers.DeleteHandler.DeleteKelas)))
}
