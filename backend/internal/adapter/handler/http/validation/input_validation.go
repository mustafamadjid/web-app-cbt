package httpx

import (
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

var disallowedSubstrings = []string{"--", "/*", "*/", "`", "$("}

func ValidateInputSafe(value, field string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	for _, item := range disallowedSubstrings {
		if strings.Contains(trimmed, item) {
			return coreerror.ErrInvalidInputSafe
		}
	}
	if strings.ContainsAny(trimmed, ";|&<>\n\r") {
		return coreerror.ErrInvalidInput
	}
	return nil
}
