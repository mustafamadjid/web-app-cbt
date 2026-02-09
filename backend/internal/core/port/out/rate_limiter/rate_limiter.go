package rate_limiter_repo

import (
	"context"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string)(allowed bool, retryAfter time.Duration, err error)
}