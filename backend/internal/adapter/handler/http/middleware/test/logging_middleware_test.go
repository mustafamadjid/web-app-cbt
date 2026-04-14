package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLogger_LogsInfoAndSetsRequestID(t *testing.T) {
	logger := &capturedLogger{}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor, ok := middleware.ActorFromContext(r.Context())
		require.True(t, ok)
		assert.Equal(t, user.ID(88), actor.IdPengguna)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{
		IdPengguna: 88,
		Role:       user.ADMIN,
		Username:   "admin",
	}))
	rec := httptest.NewRecorder()

	middleware.RequestLogger(next, logger).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Len(t, logger.withAttrs, 6)
	assert.Equal(t, "request_id", logger.withAttrs[0])
	assert.NotEmpty(t, logger.withAttrs[1])
	assert.Equal(t, "method", logger.withAttrs[2])
	assert.Equal(t, http.MethodGet, logger.withAttrs[3])
	assert.Equal(t, "path", logger.withAttrs[4])
	assert.Equal(t, "/admin/dashboard", logger.withAttrs[5])
	assert.Len(t, logger.infoCalls, 1)
	assert.Empty(t, logger.errorCalls)
	assert.Equal(t, "request completed", logger.infoCalls[0].msg)

	attrs := attrsToMap(logger.infoCalls[0].attrs)
	assert.Equal(t, "adapter.http", attrs["layer"])
	assert.Equal(t, "http_request", attrs["op"])
	assert.Equal(t, http.StatusOK, attrs["status"])
	if duration, ok := attrs["duration_ms"].(int64); ok {
		assert.GreaterOrEqual(t, duration, int64(0))
	} else {
		t.Fatalf("duration_ms has unexpected type %T", attrs["duration_ms"])
	}
	assert.Equal(t, "88", attrs["actor_id"])
	assert.Equal(t, "", attrs["err"])
	assert.NotEmpty(t, rec.Header().Get("X-Request-Id"))
}

func TestRequestLogger_LogsErrorOn5xx(t *testing.T) {
	logger := &capturedLogger{}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.RequestLogger(next, logger).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Len(t, logger.errorCalls, 1)
	assert.Empty(t, logger.infoCalls)
	assert.Equal(t, "request failed", logger.errorCalls[0].msg)

	attrs := attrsToMap(logger.errorCalls[0].attrs)
	assert.Equal(t, http.StatusInternalServerError, attrs["status"])
	assert.Equal(t, "Internal Server Error", attrs["err"])
}

func attrsToMap(attrs []any) map[string]any {
	out := make(map[string]any, len(attrs)/2)
	for i := 0; i+1 < len(attrs); i += 2 {
		key, _ := attrs[i].(string)
		out[key] = attrs[i+1]
	}
	return out
}
