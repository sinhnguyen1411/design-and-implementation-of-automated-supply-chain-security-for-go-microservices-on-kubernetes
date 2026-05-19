package token

import "testing"

func TestParse_valid(t *testing.T) {
	c, err := Parse("header.payload.sig")
	if err != nil { t.Fatal(err) }
	if c.Subject == "" { t.Fatal("empty subject") }
}

func TestParse_empty(t *testing.T) {
	_, err := Parse("")
	if err == nil { t.Fatal("expected error") }
}

func TestHasRole(t *testing.T) {
	c := Claims{Roles: []string{"trader", "admin"}}
	if !HasRole(c, "trader") { t.Fatal("should have trader") }
	if HasRole(c, "superuser") { t.Fatal("should not have superuser") }
}
