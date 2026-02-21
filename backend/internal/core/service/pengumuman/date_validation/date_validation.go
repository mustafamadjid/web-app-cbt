package pengumuman_service

import (
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

func ValidateDate(date string) error {
	t, err := time.Parse("2006-1-2", date)
	if err != nil {
		return coreerror.ErrInvalidDateFormat
	}

	if t.Format("2006-01-02") != date {
		return coreerror.ErrInvalidDateFormat
	}
	return nil
}
