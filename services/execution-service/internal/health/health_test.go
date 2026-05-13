package health

import "testing"

func TestStatus(t *testing.T) {
	if Status() == "" {
		t.Fatal("expected non-empty status")
	}
}

