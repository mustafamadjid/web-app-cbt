package routes

import "github.com/julienschmidt/httprouter"

func RegisterAuthRoutes(router *httprouter.Router, handlers AuthHandlers, mw MiddlewareContract) {
	router.POST("/auth/login", mw.RateLimiterAuth(mw.WithRequestLogger(handlers.Handler.Login)))
	router.POST("/auth/logout", mw.WithRequestLogger(mw.RateLimitStandard(handlers.Handler.Logout)))
	router.POST("/auth/refresh", mw.WithRequestLogger(mw.RateLimitStandard(handlers.Handler.Refresh)))
	router.GET("/auth/me", mw.RequireAccess(mw.RateLimitStandard(handlers.Handler.AuthMe)))
}
