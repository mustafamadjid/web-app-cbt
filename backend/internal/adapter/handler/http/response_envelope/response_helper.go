package httpx
import (
	"encoding/json"
	"net/http"
)
func WriteOK[T any](w http.ResponseWriter, status int, data T, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse[T]{
		Data:  &data,
		Message: message,
		Error: nil,
	})
}

func WriteErr(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		Data  any      `json:"data"`
		Error APIError `json:"error"`
	}{
		Data: nil,
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

