package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int           // max tokens
	tokens       int           // current tokens
	fillInterval time.Duration // time between new tokens
	lastRefill   time.Time
	mutex        sync.Mutex
}

func NewTokenBucket(capacity int, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		fillInterval: time.Second * time.Duration(10) / time.Duration(capacity),
		lastRefill:   time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	newTokens := int(elapsed / tb.fillInterval)

	if newTokens > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+newTokens)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func (tb *TokenBucket) Usage() (used int, allowed int) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	return tb.capacity - tb.tokens, tb.capacity
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
