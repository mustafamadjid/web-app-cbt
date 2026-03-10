package httpx

import (
	"regexp"
	"strings"
)

var EmailLiteRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func ValidateEmailAddress(value, field string) (string, error) {
	return ValidateByRule(strings.ToLower(value), field, Rule{
		Required:  true,
		TrimSpace: true,
		MaxLen:    254,
		RejectCtl: true,
		Pattern:   EmailLiteRe,
	})
}

func ValidateEmail(value string) (string, error) {
	return ValidateEmailAddress(value, "email")
}
