package httpx

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateRequiredPrintableText(t *testing.T) {
	value, err := ValidateRequiredPrintableText("  nama kelas  ", "nama_kelas")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "nama kelas" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidateRequiredPrintableText_Empty(t *testing.T) {
	_, err := ValidateRequiredPrintableText("   ", "nama_kelas")
	if !errors.Is(err, ErrRequiredField) {
		t.Fatalf("expected ErrRequiredField, got %v", err)
	}
}

func TestValidateUsername(t *testing.T) {
	value, err := ValidateUsername("  user.name  ")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "user.name" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidateUsername_InvalidFormat(t *testing.T) {
	_, err := ValidateUsername("user-name")
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("expected ErrInvalidFormat, got %v", err)
	}
}

func TestValidateEmailAddress(t *testing.T) {
	value, err := ValidateEmailAddress("  USER@Example.COM  ", "email")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "user@example.com" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidatePassword(t *testing.T) {
	value, err := ValidatePassword("rahasia12")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "rahasia12" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidatePassword_TooLongByBytes(t *testing.T) {
	_, err := ValidatePassword(strings.Repeat("a", PasswordMaxBytes+1))
	if !errors.Is(err, ErrTooLong) {
		t.Fatalf("expected ErrTooLong, got %v", err)
	}
}
