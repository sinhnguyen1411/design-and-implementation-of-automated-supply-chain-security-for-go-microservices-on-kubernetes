package notification

import "testing"

func TestRenderCritical(t *testing.T) {
	got := Render(Request{
		EventType: "order_rejected",
		Severity:  "critical",
		Actor:     "engine",
		Target:    "desk-a",
		Metadata:  map[string]string{"run_id": "r-001"},
	})
	if !got.Valid || got.Channel != "pagerduty" || got.Priority != "p1" {
		t.Fatalf("unexpected message: %+v", got)
	}
}

func TestRenderInvalidSeverity(t *testing.T) {
	got := Render(Request{EventType: "order_rejected", Severity: "urgent"})
	if got.Valid || got.Reason != "invalid_severity" {
		t.Fatalf("expected invalid severity, got %+v", got)
	}
}

