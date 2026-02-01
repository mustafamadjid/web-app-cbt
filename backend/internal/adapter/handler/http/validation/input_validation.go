package httpx

import (
	"fmt"
	"strings"
)

var disallowedSubstrings = []string{"--", "/*", "*/", "`", "$("}

func validateInputSafe(value, field string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	for _, item := range disallowedSubstrings {
		if strings.Contains(trimmed, item) {
			return fmt.Errorf("invalid input: %s contains invalid characters", field)
		}
	}
	if strings.ContainsAny(trimmed, ";|&<>\n\r") {
		return fmt.Errorf("invalid input: %s contains invalid characters", field)
	}
	return nil
}
