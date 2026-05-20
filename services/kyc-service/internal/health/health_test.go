package health

import "testing"

func TestStatus(t *testing.T) {
	if Status() == "" {
		t.Fatal("status should not be empty")
	}
}
