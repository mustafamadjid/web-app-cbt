package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterUserRoutes(router *httprouter.Router, userHandlers UserHandlers, resetPasswordHandlers ResetPasswordHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)

	// SISWA
	router.GET("/admin/siswa", requireAdmin(mw.RateLimitStandard(userHandlers.GetSiswaHandler.ListSiswa)))
	router.GET("/admin/siswa/:id", requireAdmin(mw.RateLimitStandard(userHandlers.GetSiswaHandler.GetSiswaByID)))
	router.POST("/admin/siswa", requireAdmin(mw.RateLimitStandard(userHandlers.CreateHandler.CreateSiswa)))
	router.PATCH("/admin/siswa/:id", requireAdmin(mw.RateLimitStandard(userHandlers.UpdateHandler.UpdateSiswa)))

	// GURU
	router.GET("/admin/guru", requireAdmin(mw.RateLimitStandard(userHandlers.GetGuruHandler.ListGuru)))
	router.GET("/admin/guru/:id", requireAdmin(mw.RateLimitStandard(userHandlers.GetGuruHandler.GetGuruByID)))
	router.POST("/admin/guru", requireAdmin(mw.RateLimitStandard(userHandlers.CreateHandler.CreateGuru)))
	router.PATCH("/admin/guru/:id", requireAdmin(mw.RateLimitStandard(userHandlers.UpdateHandler.UpdateGuru)))

	// PENGGUNA
	router.DELETE("/admin/pengguna", requireAdmin(mw.RateLimitStandard(userHandlers.DeleteHandler.DeleteUsers)))
	router.DELETE("/admin/pengguna/:id", requireAdmin(mw.RateLimitStandard(userHandlers.DeleteHandler.DeleteUser)))

	// RESET PASSWORD
	router.PUT("/admin/pengguna/:idPengguna/reset-password", requireAdmin(mw.RateLimitStandard(resetPasswordHandlers.Handler.ResetPasswordHandler)))
}
