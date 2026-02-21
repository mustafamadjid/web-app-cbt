package pengumuman_service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/date_validation"
)

func TestValidateDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		date      string
		expectErr error
	}{
		{
			name:      "branch 1 -> tanggal tidak bisa diparse",
			date:      "2026/01/01",
			expectErr: coreerror.ErrInvalidDateFormat,
		},
		{
			name:      "branch 2 -> tanggal valid parse tapi format tidak canonical",
			date:      "2026-1-2",
			expectErr: coreerror.ErrInvalidDateFormat,
		},
		{
			name:      "happy path -> format tanggal valid",
			date:      "2026-01-02",
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := pengumuman_service.ValidateDate(tt.date)
			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
