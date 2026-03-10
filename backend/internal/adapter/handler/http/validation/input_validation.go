package httpx

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrRequiredField = errors.New("required field")
	ErrTooShort      = errors.New("input too short")
	ErrTooLong       = errors.New("input too long")
	ErrInvalidFormat = errors.New("invalid format")
	ErrInvalidUTF8   = errors.New("invalid utf-8")
	ErrControlChar   = errors.New("control characters are not allowed")
)

type Rule struct {
	Required    bool
	TrimSpace   bool
	NormalizeWS bool
	MinLen      int
	MaxLen      int
	MaxBytes    int
	RejectCtl   bool
	Pattern     *regexp.Regexp
}

func ValidateByRule(value, field string, rule Rule) (string, error) {
	v := value

	if rule.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if rule.NormalizeWS {
		v = strings.Join(strings.Fields(v), " ")
	}

	if v == "" {
		if rule.Required {
			return "", fmt.Errorf("%s: %w", field, ErrRequiredField)
		}
		return "", nil
	}

	if !utf8.ValidString(v) {
		return "", fmt.Errorf("%s: %w", field, ErrInvalidUTF8)
	}

	if rule.MaxBytes > 0 && len(v) > rule.MaxBytes {
		return "", fmt.Errorf("%s: %w", field, ErrTooLong)
	}

	n := utf8.RuneCountInString(v)
	if rule.MinLen > 0 && n < rule.MinLen {
		return "", fmt.Errorf("%s: %w", field, ErrTooShort)
	}
	if rule.MaxLen > 0 && n > rule.MaxLen {
		return "", fmt.Errorf("%s: %w", field, ErrTooLong)
	}

	if rule.RejectCtl {
		for _, r := range v {
			if unicode.IsControl(r) || !unicode.IsPrint(r) {
				return "", fmt.Errorf("%s: %w", field, ErrControlChar)
			}
		}
	}

	if rule.Pattern != nil && !rule.Pattern.MatchString(v) {
		return "", fmt.Errorf("%s: %w", field, ErrInvalidFormat)
	}

	return v, nil
}

func ValidateRequiredPrintableText(value, field string) (string, error) {
	return ValidateByRule(value, field, Rule{
		Required:  true,
		TrimSpace: true,
		RejectCtl: true,
	})
}

func ValidateOptionalPrintableText(value, field string) (string, error) {
	return ValidateByRule(value, field, Rule{
		TrimSpace: true,
		RejectCtl: true,
	})
}
