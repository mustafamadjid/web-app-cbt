package tests

import (
	"testing"
	"strings"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateByRule(t *testing.T) {
	t.Parallel()

	t.Run("required printable text", func(t *testing.T) {
		value, err := httpx.ValidateRequiredPrintableText("  nama kelas  ", "nama_kelas")
		require.NoError(t, err)
		assert.Equal(t, "nama kelas", value)
	})

	t.Run("required empty", func(t *testing.T) {
		_, err := httpx.ValidateRequiredPrintableText("   ", "nama_kelas")
		assert.ErrorIs(t, err, httpx.ErrRequiredField)
	})

	t.Run("optional empty", func(t *testing.T) {
		value, err := httpx.ValidateOptionalPrintableText("   ", "catatan")
		require.NoError(t, err)
		assert.Empty(t, value)
	})

	t.Run("too short", func(t *testing.T) {
		_, err := httpx.ValidateByRule("ab", "field", httpx.Rule{MinLen: 3})
		assert.ErrorIs(t, err, httpx.ErrTooShort)
	})

	t.Run("too long by bytes", func(t *testing.T) {
		_, err := httpx.ValidateByRule(strings.Repeat("a", 5), "field", httpx.Rule{MaxBytes: 4})
		assert.ErrorIs(t, err, httpx.ErrTooLong)
	})

	t.Run("invalid utf8", func(t *testing.T) {
		_, err := httpx.ValidateByRule(string([]byte{0xff, 0xfe}), "field", httpx.Rule{})
		assert.ErrorIs(t, err, httpx.ErrInvalidUTF8)
	})

	t.Run("control char", func(t *testing.T) {
		_, err := httpx.ValidateByRule("abc\x00", "field", httpx.Rule{RejectCtl: true})
		assert.ErrorIs(t, err, httpx.ErrControlChar)
	})
}

func TestValidateUsername(t *testing.T) {
	t.Parallel()

	value, err := httpx.ValidateUsername("  user.name  ")
	require.NoError(t, err)
	assert.Equal(t, "user.name", value)

	_, err = httpx.ValidateUsername("user-name")
	assert.ErrorIs(t, err, httpx.ErrInvalidFormat)

	_, err = httpx.ValidateUsername("ab")
	assert.ErrorIs(t, err, httpx.ErrTooShort)
}

func TestValidateEmailAddress(t *testing.T) {
	t.Parallel()

	value, err := httpx.ValidateEmailAddress("  USER@Example.COM  ", "email")
	require.NoError(t, err)
	assert.Equal(t, "user@example.com", value)

	_, err = httpx.ValidateEmailAddress("not-an-email", "email")
	assert.ErrorIs(t, err, httpx.ErrInvalidFormat)
}

func TestValidateEmailAlias(t *testing.T) {
	t.Parallel()

	value, err := httpx.ValidateEmail("  Admin@Example.com  ")
	require.NoError(t, err)
	assert.Equal(t, "admin@example.com", value)
}

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	value, err := httpx.ValidatePassword("rahasia12")
	require.NoError(t, err)
	assert.Equal(t, "rahasia12", value)

	_, err = httpx.ValidatePassword("short")
	assert.ErrorIs(t, err, httpx.ErrTooShort)

	_, err = httpx.ValidatePassword(strings.Repeat("a", httpx.PasswordMaxBytes+1))
	assert.ErrorIs(t, err, httpx.ErrTooLong)
}
