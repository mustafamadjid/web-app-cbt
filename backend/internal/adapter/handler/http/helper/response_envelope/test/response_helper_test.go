package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	responsehelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseWriters(t *testing.T) {
	tests := []struct {
		name, kind    string
		status        int
		message, code string
	}{
		{name: "Branch 1 -> success writer includes typed payload", kind: "ok", status: http.StatusCreated, message: "created"},
		{name: "Branch 2 -> no-data success writer uses true sentinel", kind: "no-data", status: http.StatusOK, message: "success"},
		{name: "Branch 3 -> error writer includes structured error", kind: "error", status: http.StatusBadRequest, message: "invalid input", code: "BAD_REQUEST"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			switch tt.kind {
			case "ok":
				responsehelper.WriteOK(rec, tt.status, map[string]any{"name": "bank soal"}, tt.message)
			case "no-data":
				responsehelper.WriteOKNoData(rec, tt.status, tt.message)
			case "error":
				responsehelper.WriteErr(rec, tt.status, tt.code, tt.message)
			}
			assert.Equal(t, tt.status, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			var body map[string]any
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			if tt.kind == "ok" {
				assert.Equal(t, tt.message, body["message"])
				assert.Equal(t, "bank soal", body["data"].(map[string]any)["name"])
			}
			if tt.kind == "no-data" {
				assert.Equal(t, true, body["data"])
				assert.Equal(t, tt.message, body["message"])
			}
			if tt.kind == "error" {
				assert.Equal(t, true, body["data"])
				errBody := body["error"].(map[string]any)
				assert.Equal(t, tt.code, errBody["code"])
				assert.Equal(t, tt.message, errBody["message"])
			}
		})
	}
}
