package httpx

type APIResponse[T any] struct {
	Data *T `json:"data"`
	Message string `json:"message"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code string `json:"code"`
	Message string `json:"message"`
}