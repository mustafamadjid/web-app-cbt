package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterProfilSekolahRoutes(router *httprouter.Router, handlers ProfilSekolahHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)
	router.GET("/admin/profil-sekolah", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetProfilSekolah)))
	router.PATCH("/admin/profil-sekolah", requireAdmin(mw.RateLimitStandard(handlers.UpdateHandler.UpdateProfilSekolah)))
}
