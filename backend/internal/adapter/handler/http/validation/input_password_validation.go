package httpx

const (
	PasswordMinLength = 8
	PasswordMaxBytes  = 72
)

func ValidatePassword(value string) (string, error) {
	return ValidateByRule(value, "password", Rule{
		Required:  true,
		MinLen:    PasswordMinLength,
		MaxBytes:  PasswordMaxBytes,
		RejectCtl: true,
	})
}
