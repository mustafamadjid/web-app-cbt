package httpx

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	PersonNameRe    = regexp.MustCompile(`^[\p{L}\p{M}][\p{L}\p{M}\s.,'()-]*$`)
	GenderLabelRe   = regexp.MustCompile(`^[\p{L}\p{M}][\p{L}\p{M}\s._-]*$`)
	SafeLabelTextRe = regexp.MustCompile(`^[\p{L}\p{M}\p{N}][\p{L}\p{M}\p{N}\s.,'()&/-]*$`)
	PhoneNumberRe   = regexp.MustCompile(`^\+?[0-9][0-9\s().-]{7,19}$`)
	NIPRe           = regexp.MustCompile(`^(?:-|\d{18})$`)
	NISNRe          = regexp.MustCompile(`^(?:-|\d{10})$`)
)

func ValidatePersonName(value, field string) (string, error) {
	return ValidateByRule(value, field, Rule{
		Required:  true,
		TrimSpace: true,
		RejectCtl: true,
		Pattern:   PersonNameRe,
	})
}

func ValidateGenderLabel(value, field string) (string, error) {
	validated, err := ValidateByRule(value, field, Rule{
		Required:  true,
		TrimSpace: true,
		RejectCtl: true,
		Pattern:   GenderLabelRe,
	})
	if err != nil {
		return "", err
	}

	normalized := strings.ToUpper(validated)
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")

	switch normalized {
	case "L", "LK", "LAKI", "LAKI_LAKI", "PRIA", "MALE":
		return "LAKI_LAKI", nil
	case "P", "PR", "PEREMPUAN", "WANITA", "FEMALE":
		return "PEREMPUAN", nil
	default:
		return "", fmt.Errorf("%s: %w", field, ErrInvalidFormat)
	}
}

func ValidateSafeLabelText(value, field string) (string, error) {
	return ValidateByRule(value, field, Rule{
		Required:  true,
		TrimSpace: true,
		RejectCtl: true,
		Pattern:   SafeLabelTextRe,
	})
}

func ValidatePhoneNumber(value, field string) (string, error) {
	return ValidateByRule(value, field, Rule{
		Required:  true,
		TrimSpace: true,
		RejectCtl: true,
		Pattern:   PhoneNumberRe,
	})
}

func ValidateNIPText(value, field string) (string, error) {
	return ValidateByRule(value, field, Rule{
		Required:  true,
		TrimSpace: true,
		RejectCtl: true,
		Pattern:   NIPRe,
	})
}

func ValidateNISNText(value, field string) (string, error) {
	return ValidateByRule(value, field, Rule{
		Required:  true,
		TrimSpace: true,
		RejectCtl: true,
		Pattern:   NISNRe,
	})
}
