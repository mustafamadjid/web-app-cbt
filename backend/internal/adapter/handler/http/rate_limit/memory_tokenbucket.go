package rate_limit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter *rate.Limiter
	lastSeen time.Time
}

type MemoryTokenBucket struct {
	visitors map[string]*visitor

	rate rate.Limit
	burst int
	ttl time.Duration

	mu sync.Mutex
}

func (r *MemoryTokenBucket) cleanupLoop(ttl time.Duration) {
	ticker := time.NewTicker(ttl/2)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		for k, v := range r.visitors {
			if time.Since(v.lastSeen) > r.ttl {
				delete(r.visitors, k)
			}
		}
		r.mu.Unlock()
	}
}


func NewMemoryTokenBucket(rate rate.Limit, burst int, ttl time.Duration) *MemoryTokenBucket {
	rl := &MemoryTokenBucket{
		visitors: make(map[string]*visitor),
		rate: rate,
		burst: burst,
		ttl: ttl,
	}

	go rl.cleanupLoop(ttl)
	return rl
}

func (r *MemoryTokenBucket) Allow(ctx context.Context, key string) (allowed bool, retryAfter time.Duration, err error) {
	now := time.Now()

	r.mu.Lock()

	v,ok := r.visitors[key]
	if !ok{
		v := &visitor{
			limiter: rate.NewLimiter(r.rate, r.burst),
			lastSeen: now,
		}
		r.visitors[key] = v
	} else {
		v.lastSeen = now
	}

	res := v.limiter.ReserveN(now,1)

	r.mu.Unlock()

	if !res.OK() {
		return false, time.Second, nil
	}

	delay := res.DelayFrom(now)
	if delay > 0{
		res.CancelAt(now)
		return false, delay, nil
	}

	return true, 0, nil
}