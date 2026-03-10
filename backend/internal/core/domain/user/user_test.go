package user

import (
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
)

func TestCheckUsernameLength(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		username string
		wantErr  error
	}{
		{
			name:     "too short",
			username: "user",
			wantErr:  coreerror.ErrUsernameLengthInvalid,
		},
		{
			name:     "too long",
			username: "123456789012345678901",
			wantErr:  coreerror.ErrUsernameLengthInvalid,
		},
		{
			name:     "valid",
			username: "username123",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := CheckUsernameLength(tc.username)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}
