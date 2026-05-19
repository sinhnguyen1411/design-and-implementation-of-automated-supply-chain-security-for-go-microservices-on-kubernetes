package watchlist

type Watchlist struct {
	OwnerID string   `json:"owner_id"`
	Symbols []string `json:"symbols"`
}

func New(ownerID string) *Watchlist {
	return &Watchlist{OwnerID: ownerID, Symbols: []string{}}
}

func (w *Watchlist) Add(symbol string) bool {
	for _, s := range w.Symbols {
		if s == symbol { return false }
	}
	w.Symbols = append(w.Symbols, symbol)
	return true
}

func (w *Watchlist) Remove(symbol string) bool {
	for i, s := range w.Symbols {
		if s == symbol {
			w.Symbols = append(w.Symbols[:i], w.Symbols[i+1:]...)
			return true
		}
	}
	return false
}
