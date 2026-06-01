package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterKelasRoutes(router *httprouter.Router, handlers KelasHandlers, mw MiddlewareContract) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)

	router.GET("/admin/kelas", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.ListKelas)))
	router.GET("/admin/kelas/:idTingkatKelas/:idNamaKelas", requireAdminGuru(mw.RateLimitStandard(handlers.GetHandler.GetKelasByID)))
	router.POST("/admin/kelas/tingkat-kelas", requireAdminGuru(mw.RateLimitStandard(handlers.CreateHandler.CreateTingkatKelas)))
	router.POST("/admin/kelas/nama-kelas", requireAdminGuru(mw.RateLimitStandard(handlers.CreateHandler.CreateNamaKelas)))
	router.PATCH("/admin/kelas/nama-kelas/:idNamaKelas", requireAdminGuru(mw.RateLimitStandard(handlers.UpdateHandler.UpdateNamaKelas)))
	router.DELETE("/admin/kelas/nama-kelas/:idNamaKelas", requireAdminGuru(mw.RateLimitStandard(handlers.DeleteHandler.DeleteKelas)))
}
