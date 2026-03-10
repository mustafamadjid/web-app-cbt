package httpx

import (
	"regexp"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

var UsernameRe = regexp.MustCompile(`^[a-zA-Z0-9._]+$`)

func ValidateUsername(value string) (string, error) {
	return ValidateByRule(value, "username", Rule{
		Required:  true,
		TrimSpace: true,
		MinLen:    user.UsernameMinLength,
		MaxLen:    user.UsernameMaxLength,
		RejectCtl: true,
		Pattern:   UsernameRe,
	})
}
