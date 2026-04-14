package tests

import (
	"testing"
	"time"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatTanggalIndonesia(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, time.April, 15, 10, 30, 0, 0, time.UTC)

	assert.Equal(t, "Senin, 15 April 2024", httphelper.FormatTanggalIndonesia(ts))
}

func TestFormatTanggalWaktuIndonesia(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, time.April, 15, 10, 30, 0, 0, time.UTC)

	assert.Equal(t, "Senin 15 April 10.30", httphelper.FormatTanggalWaktuIndonesia(ts))
}

func TestFormatDateOnly(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", httphelper.FormatDateOnly(time.Time{}))
	assert.Equal(t, "2024-04-15", httphelper.FormatDateOnly(time.Date(2024, time.April, 15, 10, 30, 0, 0, time.UTC)))
}

func TestFormatTimeOnly(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "10:30", httphelper.FormatTimeOnly(time.Date(2024, time.April, 15, 10, 30, 0, 0, time.UTC)))
}

func TestFormatRFC3339(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, time.April, 15, 10, 30, 0, 0, time.UTC)
	assert.Equal(t, "2024-04-15T10:30:00Z", httphelper.FormatRFC3339(ts))
}

func TestFormatRFC3339Ptr(t *testing.T) {
	t.Parallel()

	assert.Nil(t, httphelper.FormatRFC3339Ptr(nil))

	ts := time.Date(2024, time.April, 15, 10, 30, 0, 0, time.UTC)
	value := httphelper.FormatRFC3339Ptr(&ts)
	require.NotNil(t, value)
	assert.Equal(t, "2024-04-15T10:30:00Z", *value)
}
