package account

import "testing"

func TestNew(t *testing.T) {
	a := New("owner-1")
	if a.ID == "" { t.Fatal("empty id") }
	if a.Status != KYCPending { t.Fatalf("want pending, got %s", a.Status) }
}

func TestVerify(t *testing.T) {
	a := New("owner-1")
	a.Verify()
	if a.Status != KYCVerified { t.Fatal("should be verified") }
}
