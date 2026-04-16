package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterDashboardRoutes(router *httprouter.Router, handlers DashboardHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)
	router.GET("/admin/dashboard", requireAdmin(mw.RateLimitStandard(handlers.GetHandler.GetDashboardStatistik)))
}
