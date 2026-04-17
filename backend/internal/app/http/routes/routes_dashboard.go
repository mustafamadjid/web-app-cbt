package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterDashboardRoutes(router *httprouter.Router, handlers DashboardHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)
	requireGuru := mw.RequireAccessRole(user.GURU)
	router.GET("/admin/dashboard", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetDashboardStatistik)))
	router.GET("/guru/dashboard", requireGuru(mw.RateLimitStandard(handlers.GetHandler.GetDashboardStatistik)))
}
