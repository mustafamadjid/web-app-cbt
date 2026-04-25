package integration_test

import (
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/date_validation"
	"github.com/stretchr/testify/assert"
)

func TestValidateDate(t *testing.T) {
	t.Run("format valid", func(t *testing.T) {
		t.Parallel()

		err := pengumuman_service.ValidateDate("2026-04-24")
		assert.NoError(t, err)
	})

	t.Run("format tidak valid", func(t *testing.T) {
		t.Parallel()

		err := pengumuman_service.ValidateDate("2026-4-24")
		assert.ErrorIs(t, err, coreerror.ErrInvalidDateFormat)
	})

	t.Run("tanggal tidak valid", func(t *testing.T) {
		t.Parallel()

		err := pengumuman_service.ValidateDate("2026-13-01")
		assert.ErrorIs(t, err, coreerror.ErrInvalidDateFormat)
	})
}
