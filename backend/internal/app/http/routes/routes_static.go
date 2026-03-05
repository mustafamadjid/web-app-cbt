package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func RegisterStaticRoutes(
	router *httprouter.Router,
	documentPengumumanRoute string,
	imageUploadRoute string,
	documentPengumumanHandler httprouter.Handle,
	imageUploadHandler httprouter.Handle,
	mw MiddlewareContract,
) {
	requireAdminGuru := mw.RequireAccessRole(user.ADMIN, user.GURU)

	router.GET(documentPengumumanRoute, requireAdminGuru(mw.RateLimitStandard(documentPengumumanHandler)))
	router.GET(imageUploadRoute, mw.RateLimitStandard(mw.RequireAccess(imageUploadHandler)))
}
