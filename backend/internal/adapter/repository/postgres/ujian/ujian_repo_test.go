package ujianrepo

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
)

func TestMapJadwalUjianConflictError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name: "returns conflict for jadwal exclusion violation",
			err: &pgconn.PgError{
				Code:           "23P01",
				ConstraintName: "excl_jadwal_ujian_ruangan_sesi_waktu_active",
			},
			wantErr: coreerror.ErrConflict,
		},
		{
			name: "ignores different exclusion constraint",
			err: &pgconn.PgError{
				Code:           "23P01",
				ConstraintName: "other_constraint",
			},
		},
		{
			name: "ignores non exclusion violation",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "excl_jadwal_ujian_ruangan_sesi_waktu_active",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.ErrorIs(t, mapJadwalUjianConflictError(tc.err), tc.wantErr)
		})
	}
}
