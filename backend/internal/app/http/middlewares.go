package httpmodule

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	rate_limiter_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/rate_limiter"
)

type Middlewares struct {
	logger          corelog.Logger
	accessTokenSvc  out.AccessTokenService
	refreshTokenSvc out.RefreshTokenService
	sessions        out.SessionRepository
	cookies         cookie.CookieConfig
	standardLimiter rate_limiter_repo.RateLimiter
	authLimiter     rate_limiter_repo.RateLimiter
}

func NewMiddlewares(
	logger corelog.Logger,
	accessTokenSvc out.AccessTokenService,
	refreshTokenSvc out.RefreshTokenService,
	sessions out.SessionRepository,
	cookies cookie.CookieConfig,
	standardLimiter rate_limiter_repo.RateLimiter,
	authLimiter rate_limiter_repo.RateLimiter,
) *Middlewares {
	return &Middlewares{
		logger:          logger,
		accessTokenSvc:  accessTokenSvc,
		refreshTokenSvc: refreshTokenSvc,
		sessions:        sessions,
		cookies:         cookies,
		standardLimiter: standardLimiter,
		authLimiter:     authLimiter,
	}
}

func (m *Middlewares) WithRequestLogger(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(w, r, ps)
		})
		middleware.RequestLogger(handler, m.logger).ServeHTTP(w, r)
	}
}

func (m *Middlewares) RequireAccess(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(w, r, ps)
		})
		handler = middleware.RequestLogger(handler, m.logger)
		middleware.RequireValidTokenAndSession(handler, m.accessTokenSvc, m.refreshTokenSvc, m.sessions, m.cookies).ServeHTTP(w, r)
	}
}

func (m *Middlewares) RequireAccessRole(roles ...user.Role) func(next httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next(w, r, ps)
			})
			handler = middleware.RequireActorRole(handler, roles...)
			handler = middleware.RequestLogger(handler, m.logger)
			middleware.RequireValidTokenAndSession(handler, m.accessTokenSvc, m.refreshTokenSvc, m.sessions, m.cookies).ServeHTTP(w, r)
		}
	}
}

func (m *Middlewares) RateLimiterAuth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(w, r, ps)
		})

		handler = middleware.LoginRateLimit(m.authLimiter, handler)
		handler.ServeHTTP(w, r)
	}
}

func (m *Middlewares) RateLimitStandard(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(w, r, ps)
		})
		handler = middleware.StandardRateLimit(m.standardLimiter, handler)
		handler.ServeHTTP(w, r)
	}
}
