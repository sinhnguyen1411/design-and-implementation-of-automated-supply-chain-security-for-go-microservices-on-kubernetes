package oauth2

import (
	"testing"
	"time"
)

func TestIssueAndExchangeCode(t *testing.T) {
	code := IssueCode("client-a", "https://app/callback", "user-1", []string{"read", "trade"})
	if code.Code == "" { t.Fatal("empty code") }
	if code.Used { t.Fatal("should not be used") }

	tok, err := ExchangeCode(&code, "client-a")
	if err != nil { t.Fatalf("exchange failed: %v", err) }
	if tok.Token == "" { t.Fatal("empty token") }
	if tok.ExpiresIn != 3600 { t.Fatalf("want 3600, got %d", tok.ExpiresIn) }

	_, err = ExchangeCode(&code, "client-a")
	if err == nil { t.Fatal("reuse should fail") }
}

func TestExchangeCode_clientMismatch(t *testing.T) {
	code := IssueCode("client-a", "https://app/cb", "u1", []string{"read"})
	_, err := ExchangeCode(&code, "client-b")
	if err == nil { t.Fatal("mismatched client should fail") }
}

func TestExchangeCode_expired(t *testing.T) {
	code := IssueCode("client-a", "https://app/cb", "u1", []string{"read"})
	code.ExpiresAt = time.Now().Add(-1 * time.Minute)
	_, err := ExchangeCode(&code, "client-a")
	if err == nil { t.Fatal("expired code should fail") }
}

func TestIntrospect(t *testing.T) {
	code := IssueCode("client-a", "https://app/cb", "u1", []string{"read", "trade"})
	tok, _ := ExchangeCode(&code, "client-a")
	store := map[string]AccessToken{tok.Token: tok}

	result := Introspect(tok.Token, "client-a", store)
	if !result.Active { t.Fatal("should be active") }
	if len(result.Scope) != 2 { t.Fatalf("want 2 scopes, got %d", len(result.Scope)) }

	result2 := Introspect("bad-token", "client-a", store)
	if result2.Active { t.Fatal("unknown token should be inactive") }
}
