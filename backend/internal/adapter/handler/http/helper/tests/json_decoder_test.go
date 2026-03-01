package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
)

type jsonDecoderPayload struct {
	Name string `json:"name"`
}

func TestJSONDecoder(t *testing.T) {
	t.Parallel()

	makeReq := func(body string) *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	t.Run("nil target", func(t *testing.T) {
		err := httphelper.JSONDecoder(makeReq(`{"name":"ok"}`), nil)
		assert.ErrorIs(t, err, coreerror.ErrInvalidRequestBody)
	})

	t.Run("non pointer target", func(t *testing.T) {
		var payload jsonDecoderPayload
		err := httphelper.JSONDecoder(makeReq(`{"name":"ok"}`), payload)
		assert.ErrorIs(t, err, coreerror.ErrMustBePointer)
	})

	t.Run("nil pointer target", func(t *testing.T) {
		var payload *jsonDecoderPayload
		err := httphelper.JSONDecoder(makeReq(`{"name":"ok"}`), payload)
		assert.ErrorIs(t, err, coreerror.ErrMustBePointer)
	})

	t.Run("invalid json", func(t *testing.T) {
		var payload jsonDecoderPayload
		err := httphelper.JSONDecoder(makeReq(`{"name"`), &payload)
		assert.ErrorIs(t, err, coreerror.ErrInvalidRequestBody)
	})

	t.Run("unknown field", func(t *testing.T) {
		var payload jsonDecoderPayload
		err := httphelper.JSONDecoder(makeReq(`{"name":"ok","extra":"x"}`), &payload)
		assert.ErrorIs(t, err, coreerror.ErrInvalidRequestBody)
	})

	t.Run("multiple json objects", func(t *testing.T) {
		var payload jsonDecoderPayload
		err := httphelper.JSONDecoder(makeReq(`{"name":"one"}{"name":"two"}`), &payload)
		assert.ErrorIs(t, err, coreerror.ErrInvalidRequestBody)
	})

	t.Run("success", func(t *testing.T) {
		var payload jsonDecoderPayload
		err := httphelper.JSONDecoder(makeReq(`{"name":"bank soal"}`), &payload)
		assert.NoError(t, err)
		assert.Equal(t, "bank soal", payload.Name)
	})
}
