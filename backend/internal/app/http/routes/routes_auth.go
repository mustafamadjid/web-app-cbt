package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterAuthRoutes(router *httprouter.Router, handlers AuthHandlers, mw MiddlewareContract) {
	requireAdmin := mw.RequireAccessRole(user.ADMIN)

	router.POST("/auth/login", mw.RateLimiterAuth(mw.WithRequestLogger(handlers.Handler.Login)))
	router.POST("/auth/logout", mw.WithRequestLogger(mw.RateLimitStandard(handlers.Handler.Logout)))
	router.POST("/auth/refresh", mw.WithRequestLogger(mw.RateLimitStandard(handlers.Handler.Refresh)))
	router.GET("/auth/me", mw.RequireAccess(mw.RateLimitStandard(handlers.Handler.AuthMe)))
	router.PUT("/admin/auth/revoke-session", requireAdmin(mw.RateLimitStandard(handlers.Handler.AdminRevokeUser)))
}
