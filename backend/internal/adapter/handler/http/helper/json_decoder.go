package httpx

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

func JSONDecoder(r *http.Request, v any) error {
    if v == nil {
        return coreerror.ErrInvalidRequestBody
    }
    rv := reflect.ValueOf(v)
    if rv.Kind() != reflect.Ptr || rv.IsNil() {
        return coreerror.ErrMustBePointer
    }

    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    if err := dec.Decode(v); err != nil {
        return coreerror.ErrInvalidRequestBody
    }
    if err := dec.Decode(&struct{}{}); err != io.EOF {
        return coreerror.ErrInvalidRequestBody
    }
    return nil
}
