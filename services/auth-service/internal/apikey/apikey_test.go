package apikey_test

import (
	"testing"
	"time"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/auth-service/internal/apikey"
)

func TestIssueAndValidate(t *testing.T) {
	s := apikey.NewStore()
	ttl := 24 * time.Hour
	ak, err := s.Issue("owner-1", "algo-bot", []apikey.Scope{apikey.ScopeRead, apikey.ScopeTrade}, apikey.TierAlgo, 100, &ttl)
	if err != nil {
		t.Fatal(err)
	}
	if ak.Key == "" {
		t.Fatal("key phải được trả về khi issue")
	}
	result, err := s.Validate(ak.Key, apikey.ScopeRead)
	if err != nil {
		t.Fatalf("validate thất bại: %v", err)
	}
	if result.OwnerID != "owner-1" {
		t.Fatalf("owner không đúng: %s", result.OwnerID)
	}
	if result.Key != "" {
		t.Fatal("key phải bị ẩn trong response của validate")
	}
}

func TestValidate_scopeInsufficient(t *testing.T) {
	s := apikey.NewStore()
	ak, _ := s.Issue("owner-1", "read-only", []apikey.Scope{apikey.ScopeRead}, apikey.TierFree, 10, nil)
	_, err := s.Validate(ak.Key, apikey.ScopeTrade)
	if err != apikey.ErrScopeInsufficient {
		t.Fatalf("mong đợi ErrScopeInsufficient, nhận được %v", err)
	}
}

func TestRevoke(t *testing.T) {
	s := apikey.NewStore()
	ak, _ := s.Issue("owner-1", "bot", []apikey.Scope{apikey.ScopeRead}, apikey.TierPro, 50, nil)
	if err := s.Revoke(ak.ID); err != nil {
		t.Fatal(err)
	}
	_, err := s.Validate(ak.Key, apikey.ScopeRead)
	if err != apikey.ErrKeyRevoked {
		t.Fatalf("mong đợi ErrKeyRevoked, nhận được %v", err)
	}
}

func TestRotate(t *testing.T) {
	s := apikey.NewStore()
	ak, _ := s.Issue("owner-1", "bot", []apikey.Scope{apikey.ScopeRead}, apikey.TierAlgo, 60, nil)
	oldKey := ak.Key

	newAK, err := s.Rotate(ak.ID)
	if err != nil {
		t.Fatal(err)
	}
	if newAK.Key == oldKey || newAK.Key == "" {
		t.Fatal("rotated key phải khác với key cũ và không rỗng")
	}
	if _, err = s.Validate(oldKey, apikey.ScopeRead); err == nil {
		t.Fatal("key cũ phải bị vô hiệu sau khi rotate")
	}
	if _, err = s.Validate(newAK.Key, apikey.ScopeRead); err != nil {
		t.Fatalf("key mới phải hợp lệ: %v", err)
	}
}

func TestExpiry(t *testing.T) {
	s := apikey.NewStore()
	ttl := 1 * time.Millisecond
	ak, _ := s.Issue("owner-1", "short-lived", []apikey.Scope{apikey.ScopeRead}, apikey.TierFree, 10, &ttl)
	time.Sleep(10 * time.Millisecond)
	_, err := s.Validate(ak.Key, apikey.ScopeRead)
	if err != apikey.ErrKeyExpired {
		t.Fatalf("mong đợi ErrKeyExpired, nhận được %v", err)
	}
}

func TestListByOwner(t *testing.T) {
	s := apikey.NewStore()
	s.Issue("owner-1", "key-1", []apikey.Scope{apikey.ScopeRead}, apikey.TierFree, 10, nil)
	s.Issue("owner-1", "key-2", []apikey.Scope{apikey.ScopeTrade}, apikey.TierPro, 50, nil)
	s.Issue("owner-2", "key-3", []apikey.Scope{apikey.ScopeAdmin}, apikey.TierAlgo, 100, nil)

	list := s.ListByOwner("owner-1")
	if len(list) != 2 {
		t.Fatalf("mong đợi 2 key cho owner-1, nhận được %d", len(list))
	}
	for _, item := range list {
		if item.Key != "" {
			t.Fatal("key phải bị ẩn trong danh sách")
		}
	}
}
