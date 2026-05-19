package gateway

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64
	lastFill time.Time
}

func NewTokenBucket(capacity, rate float64) *TokenBucket {
	return &TokenBucket{
		tokens:   capacity,
		capacity: capacity,
		rate:     rate,
		lastFill: time.Now(),
	}
}

func (b *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(b.lastFill).Seconds()
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.lastFill = now
}

func (b *TokenBucket) Allow() bool {
	return b.AllowN(1)
}

func (b *TokenBucket) AllowN(n float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens >= n {
		b.tokens -= n
		return true
	}
	return false
}

func (b *TokenBucket) Tokens() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	return b.tokens
}

type PerClientLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*TokenBucket
	capacity float64
	rate     float64
}

func NewPerClientLimiter(capacity, rate float64) *PerClientLimiter {
	return &PerClientLimiter{
		buckets:  make(map[string]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}
}

func (l *PerClientLimiter) Allow(clientKey string) bool {
	l.mu.Lock()
	bucket, ok := l.buckets[clientKey]
	if !ok {
		bucket = NewTokenBucket(l.capacity, l.rate)
		l.buckets[clientKey] = bucket
	}
	l.mu.Unlock()
	return bucket.Allow()
}

func (l *PerClientLimiter) Stats() map[string]float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make(map[string]float64, len(l.buckets))
	for k, b := range l.buckets {
		out[k] = b.Tokens()
	}
	return out
}
