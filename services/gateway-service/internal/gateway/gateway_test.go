package gateway

import "testing"

func TestRegistry(t *testing.T) {
	reg := NewRegistry()
	reg.Register(Route{Path: "/api/users", Method: "GET", Backend: "user-service", RateLimitRPS: 100})
	reg.Register(Route{Path: "/api/users", Method: "GET", Backend: "user-service", RateLimitRPS: 100})
	if reg.Len() != 1 { t.Fatalf("want 1, got %d", reg.Len()) }
	reg.Register(Route{Path: "/api/orders", Method: "POST", Backend: "order-service", RateLimitRPS: 50})
	if reg.Len() != 2 { t.Fatalf("want 2, got %d", reg.Len()) }
}
