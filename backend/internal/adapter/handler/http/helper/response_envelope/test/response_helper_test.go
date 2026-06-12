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

func TestWriteOK(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	payload := map[string]any{"name": "bank soal"}
	responsehelper.WriteOK(rec, http.StatusCreated, payload, "created")

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "created", body["message"])
	assert.Nil(t, body["error"])

	data, ok := body["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "bank soal", data["name"])
}

func TestWriteOKNoData(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	responsehelper.WriteOKNoData(rec, http.StatusOK, "success")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "success", body["message"])
	assert.Nil(t, body["error"])
	assert.Equal(t, true, body["data"])
}

func TestWriteErr(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	responsehelper.WriteErr(rec, http.StatusBadRequest, "BAD_REQUEST", "invalid input")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, true, body["data"])

	errObj, ok := body["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "BAD_REQUEST", errObj["code"])
	assert.Equal(t, "invalid input", errObj["message"])
}
