package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

func PreventCSRF(next http.Handler) http.Handler {
	cop := http.NewCrossOriginProtection()

	raw := os.Getenv("TRUSTED_ORIGINS")
	origins := strings.Split(raw, ",")

	for _, o := range origins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if err := cop.AddTrustedOrigin(o); err != nil {
			log.Printf("csrf: invalid trusted origin %q: %v", o, err)
		}
	}

	cop.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := corelog.FromContext(r.Context())
		logger.Error(r.Context(), "csrf origin denied",
			"layer", "adapter.http", "op", "middleware.csrf",
			"origin", r.Header.Get("Origin"),
			"referer", r.Header.Get("Referer"),
			"method", r.Method,
			"path", r.URL.Path,
		)
		http.Error(w, "CSRF check failed", http.StatusForbidden)
	}))

	return cop.Handler(next)
}
