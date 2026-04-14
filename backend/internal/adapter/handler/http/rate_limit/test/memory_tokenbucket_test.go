package tests

import (
	"context"
	"testing"
	"time"

	ratelimit "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/rate_limit"
	"golang.org/x/time/rate"
	"github.com/stretchr/testify/assert"
)

func TestMemoryTokenBucket_Allow(t *testing.T) {
	limiter := ratelimit.NewMemoryTokenBucket(rate.Every(time.Hour), 1, time.Minute)

	allowed, retryAfter, err := limiter.Allow(context.Background(), "login:alice:127.0.0.1")
	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.Zero(t, retryAfter)

	allowed, retryAfter, err = limiter.Allow(context.Background(), "login:alice:127.0.0.1")
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Greater(t, retryAfter, time.Duration(0))
}
