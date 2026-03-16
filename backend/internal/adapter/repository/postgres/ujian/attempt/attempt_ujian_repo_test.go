package attemptrepo

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
)

func TestMapAttemptUniqueViolation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name:    "active attempt constraint becomes domain error",
			err:     &pgconn.PgError{Code: "23505", ConstraintName: "uq_attempt_active"},
			wantErr: coreerror.ErrSiswaHasActiveAttempt,
		},
		{
			name:    "other unique constraint becomes conflict",
			err:     &pgconn.PgError{Code: "23505", ConstraintName: "uq_attempt_ujian_one_submitted_per_peserta"},
			wantErr: coreerror.ErrConflict,
		},
		{
			name:    "non unique pg error is ignored",
			err:     &pgconn.PgError{Code: "23503"},
			wantErr: nil,
		},
		{
			name:    "non pg error is ignored",
			err:     errors.New("boom"),
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, mapAttemptUniqueViolation(tc.err), tc.wantErr)
		})
	}
}
