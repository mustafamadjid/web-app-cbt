package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/stretchr/testify/assert"
)

func TestGetClientIP(t *testing.T) {
	t.Parallel()

	t.Run("nil request", func(t *testing.T) {
		assert.Equal(t, "", httphelper.GetClientIP(nil))
	})

	t.Run("x-forwarded-for has highest priority", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Forwarded-For", " 10.1.1.10 , 192.168.1.5 ")
		req.Header.Set("X-Real-IP", "10.1.1.99")
		req.RemoteAddr = "172.16.1.1:1234"

		assert.Equal(t, "10.1.1.10", httphelper.GetClientIP(req))
	})

	t.Run("x-real-ip fallback", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Real-IP", "10.10.10.10")
		req.RemoteAddr = "172.16.1.1:1234"

		assert.Equal(t, "10.10.10.10", httphelper.GetClientIP(req))
	})

	t.Run("remote addr host:port fallback", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "172.16.1.2:8080"

		assert.Equal(t, "172.16.1.2", httphelper.GetClientIP(req))
	})

	t.Run("remote addr raw fallback", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "172.16.1.3"

		assert.Equal(t, "172.16.1.3", httphelper.GetClientIP(req))
	})
}
