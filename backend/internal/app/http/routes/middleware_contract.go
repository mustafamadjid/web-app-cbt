package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type MiddlewareContract interface {
	WithRequestLogger(next httprouter.Handle) httprouter.Handle
	RequireAccess(next httprouter.Handle) httprouter.Handle
	RequireAccessRole(roles ...user.Role) func(next httprouter.Handle) httprouter.Handle
	RateLimiterAuth(next httprouter.Handle) httprouter.Handle
	RateLimitStandard(next httprouter.Handle) httprouter.Handle
}
