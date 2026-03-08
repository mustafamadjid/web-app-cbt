package httpx

import (
	"net/http"
	"strings"
)

func parseNullableMultipartString(req *http.Request, field string) *string {
	if req.MultipartForm == nil || req.MultipartForm.Value == nil {
		return nil
	}

	raw, ok := req.MultipartForm.Value[field]
	if !ok || len(raw) == 0 {
		return nil
	}

	val := strings.TrimSpace(raw[0])
	if val == "" || strings.EqualFold(val, "null") || strings.EqualFold(val, "nil") {
		return nil
	}

	return &val
}
