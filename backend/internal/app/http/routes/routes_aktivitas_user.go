package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterAktivitasUserRoutes(router *httprouter.Router, handlers AktivitasUserHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)
	router.GET("/admin/aktivitas-user", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetAktivitasUser)))
}
