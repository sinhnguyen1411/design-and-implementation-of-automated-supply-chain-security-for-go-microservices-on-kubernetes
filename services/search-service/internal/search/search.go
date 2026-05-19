package search

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var isinPattern = regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{9}[0-9]$`)

type TradingHours struct {
	Timezone string   `json:"timezone"`
	Open     string   `json:"open"`
	Close    string   `json:"close"`
	OpenDays []string `json:"open_days"`
}

type Symbol struct {
	Ticker       string       `json:"ticker"`
	ISIN         string       `json:"isin"`
	Name         string       `json:"name"`
	Exchange     string       `json:"exchange"`
	Sector       string       `json:"sector"`
	Currency     string       `json:"currency"`
	TradingHours TradingHours `json:"trading_hours"`
}

type Registry struct {
	byISIN   map[string]Symbol
	byTicker map[string]string
}

func NewRegistry() *Registry {
	return &Registry{
		byISIN:   make(map[string]Symbol),
		byTicker: make(map[string]string),
	}
}

var (
	ErrInvalidISIN   = errors.New("invalid ISIN format")
	ErrDuplicateISIN = errors.New("ISIN already registered")
)

func (r *Registry) Add(s Symbol) error {
	if !isinPattern.MatchString(s.ISIN) {
		return ErrInvalidISIN
	}
	if _, exists := r.byISIN[s.ISIN]; exists {
		return ErrDuplicateISIN
	}
	r.byISIN[s.ISIN] = s
	r.byTicker[strings.ToUpper(s.Ticker)] = s.ISIN
	return nil
}

func (r *Registry) ByISIN(isin string) (Symbol, bool) {
	s, ok := r.byISIN[isin]
	return s, ok
}

func (r *Registry) ByTicker(ticker string) (Symbol, bool) {
	isin, ok := r.byTicker[strings.ToUpper(ticker)]
	if !ok {
		return Symbol{}, false
	}
	return r.byISIN[isin], true
}

func (r *Registry) Search(q string) []Symbol {
	if q == "" {
		all := make([]Symbol, 0, len(r.byISIN))
		for _, s := range r.byISIN {
			all = append(all, s)
		}
		return all
	}
	q = strings.ToLower(q)
	var results []Symbol
	for _, s := range r.byISIN {
		if strings.Contains(strings.ToLower(s.Ticker), q) ||
			strings.Contains(strings.ToLower(s.Name), q) ||
			strings.Contains(strings.ToLower(s.ISIN), q) ||
			strings.Contains(strings.ToLower(s.Sector), q) {
			results = append(results, s)
		}
	}
	return results
}

func (r *Registry) IsMarketOpen(ticker string, t time.Time) (bool, error) {
	s, ok := r.ByTicker(ticker)
	if !ok {
		return false, errors.New("ticker not found: " + ticker)
	}
	loc, err := time.LoadLocation(s.TradingHours.Timezone)
	if err != nil {
		return false, err
	}
	local := t.In(loc)
	dayName := local.Weekday().String()[:3]
	dayOpen := false
	for _, d := range s.TradingHours.OpenDays {
		if d == dayName {
			dayOpen = true
			break
		}
	}
	if !dayOpen {
		return false, nil
	}
	hhmm := local.Format("15:04")
	return hhmm >= s.TradingHours.Open && hhmm < s.TradingHours.Close, nil
}
