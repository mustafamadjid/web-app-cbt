package tests

import (
	"context"
	"testing"
	"time"

	ratelimit "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/rate_limit"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestMemoryTokenBucket_Allow(t *testing.T) {
	tests := []struct {
		name        string
		setupCalls  int
		key         string
		rate        rate.Limit
		burst       int
		wantAllowed bool
		wantRetry   bool
	}{
		{name: "Branch 1 -> first request creates visitor and is allowed", key: "alice", rate: rate.Every(time.Hour), burst: 1, wantAllowed: true},
		{name: "Branch 2 -> exhausted existing visitor receives retry delay", setupCalls: 1, key: "alice", rate: rate.Every(time.Hour), burst: 1, wantRetry: true},
		{name: "Branch 3 -> zero burst reservation is not okay", key: "alice", rate: rate.Every(time.Hour), burst: 0, wantRetry: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := ratelimit.NewMemoryTokenBucket(tt.rate, tt.burst, time.Minute)
			for i := 0; i < tt.setupCalls; i++ {
				_, _, _ = limiter.Allow(context.Background(), tt.key)
			}
			allowed, retry, err := limiter.Allow(context.Background(), tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantAllowed, allowed)
			if tt.wantRetry {
				assert.Greater(t, retry, time.Duration(0))
			} else {
				assert.Zero(t, retry)
			}
		})
	}
}
