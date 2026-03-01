package httpx

import (
	"net/http"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

func MultipartHeaderBodyValidator(w http.ResponseWriter, r *http.Request, maxBytes int64) error {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		return coreerror.ErrContentTypeMustMultipart
	}

	if maxBytes <= 0 {
		maxBytes = 10 << 20
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		return coreerror.ErrInvalidMultipartForm
	}

	return nil
}
