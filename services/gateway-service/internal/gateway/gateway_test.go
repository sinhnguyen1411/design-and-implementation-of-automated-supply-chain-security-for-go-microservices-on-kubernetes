package gateway_test

import (
	"testing"
	"time"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/gateway-service/internal/gateway"
)

func TestTokenBucket_allow(t *testing.T) {
	b := gateway.NewTokenBucket(5, 10)
	for i := 0; i < 5; i++ {
		if !b.Allow() {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	if b.Allow() {
		t.Fatal("expected denial when bucket is empty")
	}
}

func TestTokenBucket_refill(t *testing.T) {
	b := gateway.NewTokenBucket(10, 100)
	for i := 0; i < 10; i++ {
		b.Allow()
	}
	if b.Allow() {
		t.Fatal("expected denial when bucket is empty")
	}
	time.Sleep(200 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected allowed after refill period")
	}
}

func TestPerClientLimiter_isolation(t *testing.T) {
	l := gateway.NewPerClientLimiter(3, 10)
	for i := 0; i < 3; i++ {
		l.Allow("client-a")
	}
	if l.Allow("client-a") {
		t.Fatal("client-a should be rate limited")
	}
	if !l.Allow("client-b") {
		t.Fatal("client-b should not be affected by client-a")
	}
}

func TestPerClientLimiter_stats(t *testing.T) {
	l := gateway.NewPerClientLimiter(5, 1)
	l.Allow("client-x")
	l.Allow("client-y")
	stats := l.Stats()
	if _, ok := stats["client-x"]; !ok {
		t.Fatal("stats should contain client-x")
	}
	if _, ok := stats["client-y"]; !ok {
		t.Fatal("stats should contain client-y")
	}
}

func TestTokenBucket_allowN(t *testing.T) {
	b := gateway.NewTokenBucket(10, 1)
	if !b.AllowN(5) {
		t.Fatal("AllowN(5) should be allowed with capacity 10")
	}
	if !b.AllowN(5) {
		t.Fatal("AllowN(5) should be allowed for remaining 5")
	}
	if b.AllowN(1) {
		t.Fatal("expected denial after exhausting bucket")
	}
}
