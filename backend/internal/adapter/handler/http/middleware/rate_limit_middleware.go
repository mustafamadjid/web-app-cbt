package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	helper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	rate_limiter_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/rate_limiter"
)

// type RateLimitMiddleware struct {
// 	limiter rate_limiter_repo.RateLimiter
// }

// func NewRateLimitMiddleware(limiter rate_limiter_repo.RateLimiter) *RateLimitMiddleware {
// 	return &RateLimitMiddleware{limiter: limiter}
// }

func StandardRateLimit(limiter rate_limiter_repo.RateLimiter,next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		actor,ok := ActorFromContext(r.Context())
		if !ok {
			next.ServeHTTP(w,r)
			return
		}
		username := actor.Username
		if username == "" {
			username = "_nouser_"
		}

		ip := helper.GetClientIP(r)
		key := "username:" + username + ":" + ip

		allowed,retryAfter,err := limiter.Allow(r.Context(),key)
		if err != nil{
			next.ServeHTTP(w,r)
			return
		}
		if !allowed {
			sec := int(retryAfter.Round(time.Second).Seconds())
			if sec < 1 {
				sec = 1
			}
			w.Header().Set("Retry-After", strconv.Itoa(sec))

			httpResponse.WriteErr(w, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "too many requests")
			return
		}

		next.ServeHTTP(w,r)
	})
}

func LoginRateLimit(limiter rate_limiter_repo.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := helper.GetClientIP(r)

		username, _, err := usernameFromJSONBody(r)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if username == "" {
			username = "_nouser_"
		}

		key := "login:" + username + ":" + ip

		allowed, retryAfter, err := limiter.Allow(r.Context(), key)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if !allowed {
			sec := int(retryAfter.Round(time.Second).Seconds())
			if sec < 1 {
				sec =1
			}
			w.Header().Set("Retry-After", strconv.Itoa(sec))
			httpResponse.WriteErr(w, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "too many requests : too many login attempts")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func usernameFromJSONBody(r *http.Request) (string, []byte, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	// restore body untuk handler berikutnya
	r.Body = io.NopCloser(bytes.NewReader(b))

	var payload struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(b, &payload); err != nil {
		return "", b, nil
	}

	u := strings.TrimSpace(payload.Username)
	if u == "" {
		return "", b, nil
	}
	return strings.ToLower(u), b, nil
}
