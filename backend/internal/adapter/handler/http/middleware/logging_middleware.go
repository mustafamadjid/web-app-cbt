package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

func RequestLogger(next http.Handler, baseLogger corelog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := newRequestID()
		logger := baseLogger.With(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
		)

		ctx := corelog.WithLogger(r.Context(), logger)
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-Id", requestID)

		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		actorID := ""
		if actor, ok := ActorFromContext(ctx); ok {
			actorID = fmt.Sprint(actor.IdPengguna)
		}

		errValue := ""
		if recorder.status >= http.StatusInternalServerError {
			errValue = http.StatusText(recorder.status)
		}

		logAttrs := []any{
			"layer", "adapter.http",
			"op", "http_request",
			"status", recorder.status,
			"duration_ms", duration.Milliseconds(),
			"actor_id", actorID,
			"err", errValue,
		}

		if recorder.status >= http.StatusInternalServerError {
			logger.Error(ctx, "request failed", logAttrs...)
			return
		}

		logger.Info(ctx, "request completed", logAttrs...)
	})
}

func newRequestID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return ""
	}
	return hex.EncodeToString(buf[:])
}
