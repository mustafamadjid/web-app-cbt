package tests

import (
	"regexp"
	"strings"
	"testing"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidateByRule(t *testing.T) {
	tests := []struct {
		name, value, field string
		rule               httpx.Rule
		want               string
		wantErr            error
	}{
		{name: "Branch 1 -> trim and normalize whitespace", value: "  hello   world  ", field: "text", rule: httpx.Rule{TrimSpace: true, NormalizeWS: true}, want: "hello world"},
		{name: "Branch 2 -> required empty value is rejected", value: "   ", field: "text", rule: httpx.Rule{Required: true, TrimSpace: true}, wantErr: httpx.ErrRequiredField},
		{name: "Branch 3 -> optional empty value is accepted", value: "   ", field: "text", rule: httpx.Rule{TrimSpace: true}},
		{name: "Branch 4 -> invalid UTF-8 is rejected", value: string([]byte{0xff}), field: "text", wantErr: httpx.ErrInvalidUTF8},
		{name: "Branch 5 -> byte limit is enforced", value: "éé", field: "text", rule: httpx.Rule{MaxBytes: 3}, wantErr: httpx.ErrTooLong},
		{name: "Branch 6 -> rune minimum is enforced", value: "ab", field: "text", rule: httpx.Rule{MinLen: 3}, wantErr: httpx.ErrTooShort},
		{name: "Branch 7 -> rune maximum is enforced", value: "abc", field: "text", rule: httpx.Rule{MaxLen: 2}, wantErr: httpx.ErrTooLong},
		{name: "Branch 8 -> control character is rejected", value: "abc\x00", field: "text", rule: httpx.Rule{RejectCtl: true}, wantErr: httpx.ErrControlChar},
		{name: "Branch 9 -> non printable rune is rejected", value: "abc\u200b", field: "text", rule: httpx.Rule{RejectCtl: true}, wantErr: httpx.ErrControlChar},
		{name: "Branch 10 -> pattern mismatch is rejected", value: "123", field: "text", rule: httpx.Rule{Pattern: regexp.MustCompile(`^[a-z]+$`)}, wantErr: httpx.ErrInvalidFormat},
		{name: "Branch 11 -> value satisfying all enabled rules is returned", value: " abc ", field: "text", rule: httpx.Rule{Required: true, TrimSpace: true, MinLen: 2, MaxLen: 4, MaxBytes: 4, RejectCtl: true, Pattern: regexp.MustCompile(`^[a-z]+$`)}, want: "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := httpx.ValidateByRule(tt.value, tt.field, tt.rule)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrintableTextValidators(t *testing.T) {
	tests := []struct {
		name, value, field, want string
		optional                 bool
		wantErr                  error
	}{
		{name: "Branch 1 -> required printable text is trimmed", value: "  nama kelas  ", field: "nama_kelas", want: "nama kelas"},
		{name: "Branch 2 -> required blank text is rejected", value: "  ", field: "nama_kelas", wantErr: httpx.ErrRequiredField},
		{name: "Branch 3 -> optional blank text is accepted", value: "  ", field: "catatan", optional: true},
		{name: "Branch 4 -> optional control character is rejected", value: "abc\x00", field: "catatan", optional: true, wantErr: httpx.ErrControlChar},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var err error
			if tt.optional {
				got, err = httpx.ValidateOptionalPrintableText(tt.value, tt.field)
			} else {
				got, err = httpx.ValidateRequiredPrintableText(tt.value, tt.field)
			}
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name, value, want string
		wantErr           error
	}{
		{name: "Path 1 -> valid username is trimmed", value: "  user.name  ", want: "user.name"},
		{name: "Path 2 -> forbidden character is rejected", value: "user-name", wantErr: httpx.ErrInvalidFormat},
		{name: "Path 3 -> short username is rejected", value: "ab", wantErr: httpx.ErrTooShort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := httpx.ValidateUsername(tt.value)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name, value, field, want string
		alias                    bool
		wantErr                  error
	}{
		{name: "Path 1 -> address is lowercased and trimmed", value: "  USER@Example.COM  ", field: "contact", want: "user@example.com"},
		{name: "Path 2 -> malformed address is rejected", value: "not-an-email", field: "contact", wantErr: httpx.ErrInvalidFormat},
		{name: "Path 3 -> ValidateEmail alias uses email field rule", value: " Admin@Example.com ", alias: true, want: "admin@example.com"},
		{name: "Path 4 -> address over maximum length is rejected", value: strings.Repeat("a", 250) + "@x.com", field: "contact", wantErr: httpx.ErrTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var err error
			if tt.alias {
				got, err = httpx.ValidateEmail(tt.value)
			} else {
				got, err = httpx.ValidateEmailAddress(tt.value, tt.field)
			}
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name, value, want string
		wantErr           error
	}{
		{name: "Path 1 -> valid password is unchanged", value: "rahasia12", want: "rahasia12"},
		{name: "Path 2 -> short password is rejected", value: "short", wantErr: httpx.ErrTooShort},
		{name: "Path 3 -> password above bcrypt byte limit is rejected", value: strings.Repeat("a", httpx.PasswordMaxBytes+1), wantErr: httpx.ErrTooLong},
		{name: "Path 4 -> password whitespace is intentionally preserved", value: " pass word ", want: " pass word "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := httpx.ValidatePassword(tt.value)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
