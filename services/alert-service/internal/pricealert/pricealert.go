package pricealert

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AlertType string

const (
	TypeAbove         AlertType = "above"
	TypeBelow         AlertType = "below"
	TypePercentChange AlertType = "percent_change"
)

type Alert struct {
	ID         string    `json:"id"`
	OwnerID    string    `json:"owner_id"`
	Symbol     string    `json:"symbol"`
	Threshold  float64   `json:"threshold"`
	Type       AlertType `json:"type"`
	Triggered  bool      `json:"triggered"`
	CreatedAt  time.Time `json:"created_at"`
	TriggeredAt *time.Time `json:"triggered_at,omitempty"`
}

type TriggerEvent struct {
	AlertID     string    `json:"alert_id"`
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	Message     string    `json:"message"`
	TriggeredAt time.Time `json:"triggered_at"`
}

type Store struct {
	mu     sync.Mutex
	alerts map[string]*Alert
}

func NewStore() *Store {
	return &Store{alerts: make(map[string]*Alert)}
}

func (s *Store) Register(ownerID, symbol string, threshold float64, t AlertType) Alert {
	a := Alert{
		ID:        uuid.NewString(),
		OwnerID:   ownerID,
		Symbol:    symbol,
		Threshold: threshold,
		Type:      t,
		CreatedAt: time.Now(),
	}
	s.mu.Lock()
	s.alerts[a.ID] = &a
	s.mu.Unlock()
	return a
}

func (s *Store) EvaluateBatch(prices map[string]float64) []TriggerEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	var events []TriggerEvent
	for _, a := range s.alerts {
		if a.Triggered {
			continue // dedup: fire only once
		}
		price, ok := prices[a.Symbol]
		if !ok {
			continue
		}
		fired, msg := evaluate(a, price)
		if fired {
			now := time.Now()
			a.Triggered = true
			a.TriggeredAt = &now
			events = append(events, TriggerEvent{
				AlertID:     a.ID,
				Symbol:      a.Symbol,
				Price:       price,
				Message:     msg,
				TriggeredAt: now,
			})
		}
	}
	return events
}

func (s *Store) List(ownerID string) []Alert {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Alert
	for _, a := range s.alerts {
		if a.OwnerID == ownerID {
			out = append(out, *a)
		}
	}
	return out
}

func (s *Store) Reset(alertID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a, ok := s.alerts[alertID]; ok {
		a.Triggered = false
		a.TriggeredAt = nil
		return true
	}
	return false
}

func evaluate(a *Alert, price float64) (bool, string) {
	switch a.Type {
	case TypeAbove:
		if price > a.Threshold {
			return true, fmt.Sprintf("%s crossed above %.4f (current: %.4f)", a.Symbol, a.Threshold, price)
		}
	case TypeBelow:
		if price < a.Threshold {
			return true, fmt.Sprintf("%s dropped below %.4f (current: %.4f)", a.Symbol, a.Threshold, price)
		}
	case TypePercentChange:
		// threshold is expected % move; price here is the % change value
		if price >= a.Threshold || price <= -a.Threshold {
			return true, fmt.Sprintf("%s moved %.2f%% (threshold: ±%.2f%%)", a.Symbol, price, a.Threshold)
		}
	}
	return false, ""
}
