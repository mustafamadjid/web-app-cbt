package httpx

import (
	"errors"
	"strings"
)

func ValidateEmail(email string) error {
	trimmedEmail := strings.TrimSpace(email)
	parts := strings.Split(trimmedEmail, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return errors.New("invalid email: email harus mengandung @ dan domain")
	}

	domain := parts[1]
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || !strings.Contains(domain, ".") {
		return errors.New("invalid email: domain email tidak valid")
	}

	return nil
}
