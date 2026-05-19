package audit_test

import (
	"encoding/json"
	"testing"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/audit-service/internal/audit"
)

func TestAppendAndReplay(t *testing.T) {
	log := audit.NewEventLog()
	payload := json.RawMessage(`{"field":"value"}`)

	e1, err := log.Append("agg-1", "user-a", audit.EventCreated, payload)
	if err != nil {
		t.Fatal(err)
	}
	e2, _ := log.Append("agg-1", "user-b", audit.EventUpdated, payload)
	log.Append("agg-2", "user-a", audit.EventCreated, payload)

	events := log.Replay("agg-1")
	if len(events) != 2 {
		t.Fatalf("expected 2 events for agg-1, got %d", len(events))
	}
	if events[0].ID != e1.ID || events[1].ID != e2.ID {
		t.Fatal("events out of order")
	}
}

func TestReplayFrom(t *testing.T) {
	log := audit.NewEventLog()
	payload := json.RawMessage(`{}`)
	log.Append("agg-1", "u", audit.EventCreated, payload)
	log.Append("agg-1", "u", audit.EventUpdated, payload)
	log.Append("agg-1", "u", audit.EventApproved, payload)

	events := log.ReplayFrom("agg-1", 1)
	if len(events) != 2 {
		t.Fatalf("expected 2 events from version 1, got %d", len(events))
	}
}

func TestSnapshot(t *testing.T) {
	log := audit.NewEventLog()
	payload := json.RawMessage(`{}`)
	log.Append("agg-1", "u", audit.EventCreated, payload)
	log.Append("agg-1", "u", audit.EventUpdated, payload)

	state := json.RawMessage(`{"count":2}`)
	snap := log.TakeSnapshot("agg-1", state)
	if snap.Version != 2 {
		t.Fatalf("expected snapshot version 2, got %d", snap.Version)
	}

	retrieved, ok := log.LatestSnapshot("agg-1")
	if !ok {
		t.Fatal("snapshot not found")
	}
	if retrieved.Version != snap.Version {
		t.Fatal("snapshot version mismatch")
	}
}

func TestAppend_emptyAggregate(t *testing.T) {
	log := audit.NewEventLog()
	_, err := log.Append("", "u", audit.EventCreated, nil)
	if err == nil {
		t.Fatal("expected error for empty aggregate_id")
	}
}

func TestVersionMonotone(t *testing.T) {
	log := audit.NewEventLog()
	payload := json.RawMessage(`{}`)
	var versions []int
	for i := 0; i < 5; i++ {
		e, _ := log.Append("agg-1", "u", audit.EventUpdated, payload)
		versions = append(versions, e.Version)
	}
	for i := 1; i < len(versions); i++ {
		if versions[i] != versions[i-1]+1 {
			t.Fatalf("versions not monotone: %v", versions)
		}
	}
}
