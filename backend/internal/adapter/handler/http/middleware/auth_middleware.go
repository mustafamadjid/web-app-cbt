package middleware

import (
	"context"
	"net/http"
	"time"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
)

type ctxKey string

const actorKey ctxKey = "actor"

func ActorFromContext(ctx context.Context) (user.Actor, bool) {
	actor,ok := ctx.Value(actorKey).(user.Actor)
	return actor, ok
}

func WithActor(next http.Handler, verifier out.AccessTokenService, cookieName string) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, req *http.Request) {
		c,err := req.Cookie(cookieName)
		if err != nil || c.Value == "" {
			next.ServeHTTP(write, req)
			return
		}

		idPengguna, role, username, err := verifier.VerifyAccessToken(c.Value,time.Now())
		if err != nil {
			next.ServeHTTP(write, req)
			return
		}

		actor := user.Actor {
			IdPengguna: idPengguna,
			Role: role,
			Username: username,
		}
		ctx := context.WithValue(req.Context(), actorKey, actor)
		next.ServeHTTP(write, req.WithContext(ctx))
	})
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, req *http.Request) {
		_, ok := ActorFromContext(req.Context())
		if !ok {
			httpResponse.WriteErr(write,http.StatusUnauthorized,"UNAUTHORIZED","unauthorized")
			return
		}
		next.ServeHTTP(write,req)
	})
}



