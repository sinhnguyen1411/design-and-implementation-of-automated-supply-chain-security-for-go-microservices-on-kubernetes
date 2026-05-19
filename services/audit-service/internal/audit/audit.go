package audit

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventCreated  EventType = "created"
	EventUpdated  EventType = "updated"
	EventDeleted  EventType = "deleted"
	EventApproved EventType = "approved"
	EventRejected EventType = "rejected"
)

type Event struct {
	ID          string          `json:"id"`
	AggregateID string          `json:"aggregate_id"`
	Type        EventType       `json:"type"`
	Actor       string          `json:"actor"`
	Payload     json.RawMessage `json:"payload"`
	Version     int             `json:"version"`
	CreatedAt   time.Time       `json:"created_at"`
}

type Snapshot struct {
	AggregateID string          `json:"aggregate_id"`
	Version     int             `json:"version"`
	State       json.RawMessage `json:"state"`
	TakenAt     time.Time       `json:"taken_at"`
}

type EventLog struct {
	mu        sync.RWMutex
	events    []Event
	snapshots map[string]Snapshot
	version   map[string]int
}

func NewEventLog() *EventLog {
	return &EventLog{
		snapshots: make(map[string]Snapshot),
		version:   make(map[string]int),
	}
}

var ErrInvalidAggregate = errors.New("aggregate_id must not be empty")

func (l *EventLog) Append(aggregateID, actor string, t EventType, payload json.RawMessage) (Event, error) {
	if aggregateID == "" {
		return Event{}, ErrInvalidAggregate
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.version[aggregateID]++
	e := Event{
		ID:          uuid.NewString(),
		AggregateID: aggregateID,
		Type:        t,
		Actor:       actor,
		Payload:     payload,
		Version:     l.version[aggregateID],
		CreatedAt:   time.Now().UTC(),
	}
	l.events = append(l.events, e)
	return e, nil
}

func (l *EventLog) Replay(aggregateID string) []Event {
	return l.ReplayFrom(aggregateID, 0)
}

func (l *EventLog) ReplayFrom(aggregateID string, fromVersion int) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var result []Event
	for _, e := range l.events {
		if e.AggregateID == aggregateID && e.Version > fromVersion {
			result = append(result, e)
		}
	}
	return result
}

func (l *EventLog) TakeSnapshot(aggregateID string, state json.RawMessage) Snapshot {
	l.mu.Lock()
	defer l.mu.Unlock()
	snap := Snapshot{
		AggregateID: aggregateID,
		Version:     l.version[aggregateID],
		State:       state,
		TakenAt:     time.Now().UTC(),
	}
	l.snapshots[aggregateID] = snap
	return snap
}

func (l *EventLog) LatestSnapshot(aggregateID string) (Snapshot, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	snap, ok := l.snapshots[aggregateID]
	return snap, ok
}

func (l *EventLog) All() []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}
