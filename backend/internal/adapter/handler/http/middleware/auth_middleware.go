package middleware

import (
	"context"
	"net/http"
	"time"

	// "github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
)

type ctxKey string

const actorKey ctxKey = "actor"

func ActorFromContext(ctx context.Context) (user.Actor, bool) {
	actor, ok := ctx.Value(actorKey).(user.Actor)
	return actor, ok
}

func RequireValidAccessToken(next http.Handler, access out.AccessTokenService, cookies cookie.CookieConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookies.AccessName)
		if err != nil || c.Value == "" {
			httpResponse.WriteErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
			return
		}

		uid, role, username, err := access.VerifyAccessToken(c.Value, time.Now())
		if err != nil {
			cookie.ClearAuthCookies(w, cookies)
			httpResponse.WriteErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
			return
		}

		actor := user.Actor{
			IdPengguna: uid,
			Role:       role,
			Username:   username,
		}

		ctx := context.WithValue(r.Context(), actorKey, actor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireActorRole(next http.Handler, allowedRoles ...user.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor, ok := ActorFromContext(r.Context())
		if !ok {
			httpResponse.WriteErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
			return
		}

		for _, role := range allowedRoles {
			if actor.Role == role {
				next.ServeHTTP(w, r)
				return
			}
		}

		httpResponse.WriteErr(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
	})
}

// func RequireAccess(next httprouter.Handle, tokens out.AccessTokenService, cookies cookie.CookieConfig) httprouter.Handle{
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			next(w, r, ps)
// 		})
// 		RequireValidAccessToken(handler, tokens, cookies).ServeHTTP(w, r)
// 	}
// }
